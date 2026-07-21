package android

import (
	"fmt"
	"strings"

	"github.com/Baba01hacker666/Go2APK/internal/config"
	"github.com/Baba01hacker666/Go2APK/ir"
)

// RenderNativeBridge creates the Java JNI bridge class with all required methods.
func RenderNativeBridge(cfg config.Config) string {
	return fmt.Sprintf(`package %s;

import android.app.Activity;
import android.content.Intent;
import android.net.Uri;
import android.Manifest;
import android.content.pm.PackageManager;

final class NativeBridge {
    private static final String LIBRARY_NAME = "go2apkapp";

    static {
        System.loadLibrary(LIBRARY_NAME);
    }

    private NativeBridge() {}

    // ── JNI entry points (called from Go) ────────────────────────────────────

    private static native void sendEventToGo(String eventName);
    private static native void onPermissionResult(String permission, boolean granted);
    static native void onVpnEstablished(int fd);

    // ── Java→Go dispatch ─────────────────────────────────────────────────────

    static void sendEvent(String eventName) {
        sendEventToGo(eventName);
    }

    // ── Activity & Context helpers ────────────────────────────────────────────

    private static Activity currentActivity;

    public static void setActivity(Activity activity) {
        currentActivity = activity;
    }

    // ── UI Updates ────────────────────────────────────────────────────────────

    /** Called from Go to update any widget text on the UI thread. */
    public static void updateText(String id, String text) {
        if (currentActivity instanceof Go2ApkActivity) {
            currentActivity.runOnUiThread(() ->
                ((Go2ApkActivity) currentActivity).updateWidgetText(id, text));
        }
    }

    /** Called from Go to animate a widget on the UI thread. */
    public static void animate(String id, String property, float to, int durationMs) {
        if (currentActivity instanceof Go2ApkActivity) {
            currentActivity.runOnUiThread(() ->
                ((Go2ApkActivity) currentActivity).animateWidget(id, property, to, durationMs));
        }
    }

    /** Called from Go to read the current text of a widget. */
    public static String getText(String id) {
        if (currentActivity instanceof Go2ApkActivity) {
            return ((Go2ApkActivity) currentActivity).getWidgetText(id);
        }
        return "";
    }

    // ── Intents ───────────────────────────────────────────────────────────────

    /**
     * Called from Go to start an Android activity.
     * @param action  Intent action string (e.g. "android.intent.action.VIEW")
     * @param data    URI data string, or "" if none
     * @param pkg     explicit package name for the target app, or "" if implicit
     */
    public static void startActivity(String action, String data, String pkg) {
        if (currentActivity == null) return;
        Intent intent = new Intent(action);
        if (!data.isEmpty()) intent.setData(Uri.parse(data));
        if (!pkg.isEmpty())  intent.setPackage(pkg);
        currentActivity.startActivity(intent);
    }

    // ── Broadcasts ────────────────────────────────────────────────────────────

    /** Called from Go to send a local broadcast. */
    public static void sendBroadcast(String action) {
        if (currentActivity == null) return;
        Intent intent = new Intent(action);
        currentActivity.sendBroadcast(intent);
    }

    // ── Permissions ───────────────────────────────────────────────────────────

    private static final int PERMISSION_REQUEST_CODE = 1001;

    /**
     * Called from Go to request a dangerous runtime permission.
     * The result is delivered back via onPermissionResult on the Go side.
     */
    public static void requestPermission(String permission) {
        if (currentActivity == null) return;
        if (currentActivity.checkSelfPermission(permission) == PackageManager.PERMISSION_GRANTED) {
            onPermissionResult(permission, true);
            return;
        }
        currentActivity.requestPermissions(new String[]{permission}, PERMISSION_REQUEST_CODE);
    }

    /** Called by MainActivity.onRequestPermissionsResult to relay the result to Go. */
    static void deliverPermissionResult(String permission, boolean granted) {
        onPermissionResult(permission, granted);
    }

    // ── VPN ───────────────────────────────────────────────────────────────────

    private static String pendingVpnConfig = null;
    private static final int VPN_REQUEST_CODE = 1002;

    public static void startVpn(String configJson) {
        if (currentActivity == null) return;
        Intent intent = android.net.VpnService.prepare(currentActivity);
        if (intent != null) {
            pendingVpnConfig = configJson;
            currentActivity.startActivityForResult(intent, VPN_REQUEST_CODE);
        } else {
            // Already authorized
            launchVpnService(configJson);
        }
    }

    static void handleActivityResult(int requestCode, int resultCode, Intent data) {
        if (requestCode == VPN_REQUEST_CODE && resultCode == Activity.RESULT_OK) {
            launchVpnService(pendingVpnConfig);
            pendingVpnConfig = null;
        }
    }

    public static void navigate(String target) {
        if (currentActivity == null) return;
        try {
            Class<?> targetClass = Class.forName(currentActivity.getPackageName() + "." + target);
            Intent intent = new Intent(currentActivity, targetClass);
            currentActivity.startActivity(intent);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private static void launchVpnService(String configJson) {
        if (currentActivity == null) return;
        try {
            Class<?> vpnClass = Class.forName(currentActivity.getPackageName() + ".Go2ApkVpnService");
            Intent intent = new Intent(currentActivity, vpnClass);
            intent.putExtra("config", configJson);
            currentActivity.startService(intent);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
`, cfg.Package)
}

// RenderBroadcastReceiver generates a Java BroadcastReceiver that routes
// received intents back to Go through NativeBridge.sendEvent.
func RenderBroadcastReceiver(cfg config.Config, receivers []ir.BroadcastReceiverDecl) string {
	if len(receivers) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("package %s;\n\n", cfg.Package))
	sb.WriteString("import android.content.BroadcastReceiver;\n")
	sb.WriteString("import android.content.Context;\n")
	sb.WriteString("import android.content.Intent;\n\n")
	sb.WriteString("public class Go2APKBroadcastReceiver extends BroadcastReceiver {\n")
	sb.WriteString("    @Override\n")
	sb.WriteString("    public void onReceive(Context context, Intent intent) {\n")
	sb.WriteString("        String action = intent.getAction();\n")
	sb.WriteString("        if (action == null) return;\n")
	for _, recv := range receivers {
		sb.WriteString(fmt.Sprintf("        if (action.equals(%q)) { NativeBridge.sendEvent(%q); return; }\n",
			recv.Action, recv.Name))
	}
	sb.WriteString("    }\n")
	sb.WriteString("}\n")
	return sb.String()
}

// RenderVpnService generates the Java VpnService class if VPN is requested.
func RenderVpnService(cfg config.Config) string {
	return fmt.Sprintf(`package %s;

import android.content.Intent;
import android.net.VpnService;
import android.os.ParcelFileDescriptor;
import org.json.JSONObject;

public class Go2ApkVpnService extends VpnService {
    private ParcelFileDescriptor mInterface;

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        if (intent != null && "STOP".equals(intent.getAction())) {
            stopSelf();
            return START_NOT_STICKY;
        }

        String configJson = intent.getStringExtra("config");
        if (configJson == null) return START_NOT_STICKY;

        Builder builder = new Builder();
        try {
            JSONObject config = new JSONObject(configJson);
            if (config.has("address")) {
                builder.addAddress(config.getString("address"), config.optInt("prefix", 24));
            }
            if (config.has("route")) {
                builder.addRoute(config.getString("route"), config.optInt("routePrefix", 0));
            }
            if (config.has("dns")) {
                builder.addDnsServer(config.getString("dns"));
            }
            if (config.has("mtu")) {
                builder.setMtu(config.getInt("mtu"));
            }
            if (config.has("session")) {
                builder.setSession(config.getString("session"));
            }

            mInterface = builder.establish();
            if (mInterface != null) {
                int fd = mInterface.getFd();
                NativeBridge.onVpnEstablished(fd);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }

        return START_STICKY;
    }

    @Override
    public void onDestroy() {
        try {
            if (mInterface != null) mInterface.close();
        } catch (Exception e) {}
        super.onDestroy();
    }
}
`, cfg.Package)
}

func RenderGo2ApkActivity(cfg config.Config) string {
	return fmt.Sprintf(`package %s;

public interface Go2ApkActivity {
    void updateWidgetText(String id, String text);
    void animateWidget(String id, String property, float to, int durationMs);
    String getWidgetText(String id);
}
`, cfg.Package)
}

// RenderDynamicActivity generates an Android Activity class for a given UI widget tree.
func RenderDynamicActivity(cfg config.Config, activityName string, root ir.Widget) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(`package %s;

import android.content.Intent;
import android.app.Activity;
import android.graphics.Color;
import android.os.Bundle;
import android.view.Gravity;
import android.view.View;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.VideoView;
import android.media.MediaPlayer;
import android.net.Uri;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import java.net.URL;
import java.util.concurrent.Executors;
import android.animation.ObjectAnimator;

public class %s extends Activity implements Go2ApkActivity {
`, cfg.Package, activityName))

	if root != nil {
		declareFields(&builder, root)
	}

	builder.WriteString(fmt.Sprintf(`
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("%s");
        setContentView(createAppView());
        NativeBridge.setActivity(this);
    }

    private View createAppView() {
`, cfg.Name))

	if root != nil {
		counter := 0
		rootView := buildView(&builder, root, "", &counter)
		builder.WriteString(fmt.Sprintf("        return %s;\n", rootView))
	} else {
		// Fallback blank view
		builder.WriteString("        LinearLayout rootView = new LinearLayout(this);\n")
		builder.WriteString("        return rootView;\n")
	}

	builder.WriteString(`    }

    public void updateWidgetText(String id, String text) {
`)
	if root != nil {
		writeUpdateWidgetCases(&builder, root)
	}
	builder.WriteString(`    }

    public void animateWidget(String id, String property, float to, int durationMs) {
`)
	if root != nil {
		writeAnimateWidgetCases(&builder, root)
	}
	builder.WriteString(`    }

    public String getWidgetText(String id) {
`)
	if root != nil {
		writeGetWidgetCases(&builder, root)
	}
	builder.WriteString(`        return "";
    }

    private void loadImage(ImageView view, String src) {
        if (src.startsWith("http")) {
            Executors.newSingleThreadExecutor().execute(() -> {
                try {
                    Bitmap bmp = BitmapFactory.decodeStream(new URL(src).openConnection().getInputStream());
                    runOnUiThread(() -> view.setImageBitmap(bmp));
                } catch (Exception e) { e.printStackTrace(); }
            });
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        NativeBridge.handleActivityResult(requestCode, resultCode, data);
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, String[] permissions, int[] grantResults) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults);
        for (int i = 0; i < permissions.length; i++) {
            boolean granted = (grantResults[i] == android.content.pm.PackageManager.PERMISSION_GRANTED);
            NativeBridge.deliverPermissionResult(permissions[i], granted);
        }
    }
}
`)
	return builder.String()
}

func writeUpdateWidgetCases(b *strings.Builder, w ir.Widget) {
	switch v := w.(type) {
	case ir.ColumnWidget:
		for _, child := range v.Children {
			writeUpdateWidgetCases(b, child)
		}
	case ir.RowWidget:
		for _, child := range v.Children {
			writeUpdateWidgetCases(b, child)
		}
	case ir.TextViewWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) this.%s.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }\n", v.ID, v.ID, v.ID))
		}
	case ir.ButtonWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) this.%s.setText(android.text.Html.fromHtml(text, android.text.Html.FROM_HTML_MODE_COMPACT)); return; }\n", v.ID, v.ID, v.ID))
		}
	case ir.TextFieldWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) this.%s.setText(text); return; }\n", v.ID, v.ID, v.ID))
		}
	}
}

func writeAnimateWidgetCases(b *strings.Builder, w ir.Widget) {
	switch v := w.(type) {
	case ir.ColumnWidget:
		for _, child := range v.Children {
			writeAnimateWidgetCases(b, child)
		}
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) ObjectAnimator.ofFloat(this.%s, property, to).setDuration(durationMs).start(); return; }\n", v.ID, v.ID, v.ID))
		}
	case ir.RowWidget:
		for _, child := range v.Children {
			writeAnimateWidgetCases(b, child)
		}
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) ObjectAnimator.ofFloat(this.%s, property, to).setDuration(durationMs).start(); return; }\n", v.ID, v.ID, v.ID))
		}
	case ir.TextViewWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) ObjectAnimator.ofFloat(this.%s, property, to).setDuration(durationMs).start(); return; }\n", v.ID, v.ID, v.ID))
		}
	case ir.ButtonWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) ObjectAnimator.ofFloat(this.%s, property, to).setDuration(durationMs).start(); return; }\n", v.ID, v.ID, v.ID))
		}
	case ir.TextFieldWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) ObjectAnimator.ofFloat(this.%s, property, to).setDuration(durationMs).start(); return; }\n", v.ID, v.ID, v.ID))
		}
	case ir.ImageWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) ObjectAnimator.ofFloat(this.%s, property, to).setDuration(durationMs).start(); return; }\n", v.ID, v.ID, v.ID))
		}
	case ir.VideoWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) ObjectAnimator.ofFloat(this.%s, property, to).setDuration(durationMs).start(); return; }\n", v.ID, v.ID, v.ID))
		}
	}
}

func writeGetWidgetCases(b *strings.Builder, w ir.Widget) {
	switch v := w.(type) {
	case ir.ColumnWidget:
		for _, child := range v.Children {
			writeGetWidgetCases(b, child)
		}
	case ir.RowWidget:
		for _, child := range v.Children {
			writeGetWidgetCases(b, child)
		}
	case ir.TextViewWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\") && this.%s != null) return this.%s.getText().toString();\n", v.ID, v.ID, v.ID))
		}
	case ir.TextFieldWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\") && this.%s != null) return this.%s.getText().toString();\n", v.ID, v.ID, v.ID))
		}
	}
}

func declareFields(b *strings.Builder, w ir.Widget) {
	switch v := w.(type) {
	case ir.ColumnWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("    private LinearLayout %s;\n", v.ID))
		}
		for _, child := range v.Children {
			declareFields(b, child)
		}
	case ir.RowWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("    private LinearLayout %s;\n", v.ID))
		}
		for _, child := range v.Children {
			declareFields(b, child)
		}
	case ir.TextViewWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("    private TextView %s;\n", v.ID))
		}
	case ir.ButtonWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("    private Button %s;\n", v.ID))
		}
	case ir.TextFieldWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("    private EditText %s;\n", v.ID))
		}
	case ir.ImageWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("    private ImageView %s;\n", v.ID))
		}
	case ir.AudioWidget:
		// AudioWidget has no UI, but we could declare it as MediaPlayer if needed, though we don't have to unless we need to animate it, which doesn't make sense. But for consistency:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("    private MediaPlayer %s;\n", v.ID))
		}
	case ir.VideoWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("    private VideoView %s;\n", v.ID))
		}
	case ir.XMLView:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("    private View %s;\n", v.ID))
		}
	}
}

func applyStyle(b *strings.Builder, viewVar string, style ir.Style, defaultWidth, defaultHeight int, defaultWeight float32) {
	width := defaultWidth
	if style.Width != 0 {
		width = style.Width
	}
	height := defaultHeight
	if style.Height != 0 {
		height = style.Height
	}
	weight := defaultWeight
	if style.Weight != 0 {
		weight = style.Weight
	}

	widthStr := fmt.Sprintf("%d", width)
	if width == -1 {
		widthStr = "LinearLayout.LayoutParams.MATCH_PARENT"
	} else if width == -2 {
		widthStr = "LinearLayout.LayoutParams.WRAP_CONTENT"
	} else if width == 0 && weight > 0 {
		widthStr = "0"
	}

	heightStr := fmt.Sprintf("%d", height)
	if height == -1 {
		heightStr = "LinearLayout.LayoutParams.MATCH_PARENT"
	} else if height == -2 {
		heightStr = "LinearLayout.LayoutParams.WRAP_CONTENT"
	}

	// If weight is set, Android prefers the scaling dimension to be 0
	if weight > 0 {
		// Heuristic: if we're a Row (usually width is MATCH, height is WRAP), and we want it to scale vertically, we should set height to 0.
		// Wait, we don't know the parent's orientation here easily. We'll just let Android handle WRAP_CONTENT with weight (it works, but 0 is more efficient).
		// Actually, if it's a Button in a horizontal Row, its width is usually set to 0.
		if width == 0 {
			widthStr = "0"
		}
		if height == 0 {
			heightStr = "0"
		}
	}

	b.WriteString(fmt.Sprintf("        LinearLayout.LayoutParams lp_%s = new LinearLayout.LayoutParams(%s, %s, %ff);\n", viewVar, widthStr, heightStr, weight))
	if style.Margin != 0 {
		b.WriteString(fmt.Sprintf("        lp_%s.setMargins(%d, %d, %d, %d);\n", viewVar, style.Margin, style.Margin, style.Margin, style.Margin))
	} else {
		// Default margins if we don't specify, we'll keep 0
	}
	b.WriteString(fmt.Sprintf("        %s.setLayoutParams(lp_%s);\n", viewVar, viewVar))

	if style.Padding != 0 {
		b.WriteString(fmt.Sprintf("        %s.setPadding(%d, %d, %d, %d);\n", viewVar, style.Padding, style.Padding, style.Padding, style.Padding))
	}
	if style.BackgroundColor != "" {
		b.WriteString(fmt.Sprintf("        %s.setBackgroundColor(Color.parseColor(\"%s\"));\n", viewVar, style.BackgroundColor))
	}

	if style.Animation.Type != "" {
		duration := style.Animation.Duration
		if duration == 0 {
			duration = 300 // default duration
		}
		delay := style.Animation.Delay
		b.WriteString(fmt.Sprintf("        %s.postDelayed(() -> {\n", viewVar))
		switch style.Animation.Type {
		case "fade_in":
			b.WriteString(fmt.Sprintf("            %s.setAlpha(0f);\n", viewVar))
			b.WriteString(fmt.Sprintf("            %s.animate().alpha(1f).setDuration(%d).start();\n", viewVar, duration))
		case "slide_up":
			b.WriteString(fmt.Sprintf("            %s.setTranslationY(100f);\n", viewVar))
			b.WriteString(fmt.Sprintf("            %s.setAlpha(0f);\n", viewVar))
			b.WriteString(fmt.Sprintf("            %s.animate().translationY(0f).alpha(1f).setDuration(%d).start();\n", viewVar, duration))
		case "bounce":
			b.WriteString(fmt.Sprintf("            %s.animate().translationYBy(-50f).setDuration(%d/2).withEndAction(() -> %s.animate().translationYBy(50f).setDuration(%d/2).start()).start();\n", viewVar, duration, viewVar, duration))
		case "pulse":
			b.WriteString(fmt.Sprintf("            %s.animate().scaleX(1.1f).scaleY(1.1f).setDuration(%d/2).withEndAction(() -> %s.animate().scaleX(1f).scaleY(1f).setDuration(%d/2).start()).start();\n", viewVar, duration, viewVar, duration))
		}
		b.WriteString(fmt.Sprintf("        }, %d);\n", delay))
	}
}

func applyTextStyle(b *strings.Builder, viewVar string, style ir.Style, defaultSize int) {
	size := defaultSize
	if style.TextSize != 0 {
		size = style.TextSize
	}
	b.WriteString(fmt.Sprintf("        %s.setTextSize(%d);\n", viewVar, size))
	if style.TextColor != "" {
		b.WriteString(fmt.Sprintf("        %s.setTextColor(Color.parseColor(\"%s\"));\n", viewVar, style.TextColor))
	}
}

func applyCSS(b *strings.Builder, viewVar string, css string) {
	if strings.TrimSpace(css) == "" {
		return
	}
	rules := strings.Split(css, ";")

	hasBackground := false
	bgColor := ""
	hasRadius := false
	radius := ""
	hasBorder := false
	borderWidth := ""
	borderColor := ""

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			continue
		}
		prop := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch prop {
		case "background-color":
			hasBackground = true
			bgColor = val
		case "border-radius":
			hasRadius = true
			radius = strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(val, "px"), "dp"), "sp")
		case "border":
			// format: "1px solid black" or similar
			parts := strings.Split(val, " ")
			if len(parts) >= 3 {
				hasBorder = true
				borderWidth = strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(parts[0], "px"), "dp"), "sp")
				borderColor = parts[2]
			}
		case "box-shadow", "elevation":
			// We only parse elevation for Android natively
			val = strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(val, "px"), "dp"), "sp")
			b.WriteString(fmt.Sprintf("        %s.setElevation(Float.parseFloat(\"%s\"));\n", viewVar, val))
		case "color":
			b.WriteString(fmt.Sprintf("        if (%s instanceof TextView) ((TextView)%s).setTextColor(Color.parseColor(\"%s\"));\n", viewVar, viewVar, val))
		case "font-size":
			size := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(val, "px"), "dp"), "sp")
			b.WriteString(fmt.Sprintf("        if (%s instanceof TextView) ((TextView)%s).setTextSize(Float.parseFloat(\"%s\"));\n", viewVar, viewVar, size))
		case "padding":
			pad := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(val, "px"), "dp"), "sp")
			b.WriteString(fmt.Sprintf("        %s.setPadding((int)Float.parseFloat(\"%s\"), (int)Float.parseFloat(\"%s\"), (int)Float.parseFloat(\"%s\"), (int)Float.parseFloat(\"%s\"));\n", viewVar, pad, pad, pad, pad))
		case "text-align":
			var gravity string
			switch val {
			case "center":
				gravity = "Gravity.CENTER"
			case "right", "end":
				gravity = "Gravity.END"
			case "left", "start":
				gravity = "Gravity.START"
			}
			if gravity != "" {
				b.WriteString(fmt.Sprintf("        if (%s instanceof TextView) ((TextView)%s).setGravity(%s);\n", viewVar, viewVar, gravity))
			}
		case "font-weight":
			if val == "bold" {
				b.WriteString(fmt.Sprintf("        if (%s instanceof TextView) ((TextView)%s).setTypeface(android.graphics.Typeface.DEFAULT_BOLD);\n", viewVar, viewVar))
			}
		case "opacity":
			b.WriteString(fmt.Sprintf("        %s.setAlpha(Float.parseFloat(\"%s\"));\n", viewVar, val))
		}
	}

	// Apply background drawable if needed
	if hasBackground || hasRadius || hasBorder {
		b.WriteString("        {\n")
		b.WriteString("            android.graphics.drawable.GradientDrawable gd = new android.graphics.drawable.GradientDrawable();\n")
		if hasBackground {
			b.WriteString(fmt.Sprintf("            gd.setColor(Color.parseColor(\"%s\"));\n", bgColor))
		}
		if hasRadius {
			b.WriteString(fmt.Sprintf("            gd.setCornerRadius(Float.parseFloat(\"%s\"));\n", radius))
		}
		if hasBorder {
			b.WriteString(fmt.Sprintf("            gd.setStroke((int)Float.parseFloat(\"%s\"), Color.parseColor(\"%s\"));\n", borderWidth, borderColor))
		}
		b.WriteString(fmt.Sprintf("            %s.setBackground(gd);\n", viewVar))
		b.WriteString("        }\n")
	}
}

func buildView(b *strings.Builder, w ir.Widget, parentVar string, counter *int) string {
	*counter++
	viewVar := fmt.Sprintf("view%d", *counter)

	switch v := w.(type) {
	case ir.ColumnWidget:
		b.WriteString(fmt.Sprintf("        LinearLayout %s = new LinearLayout(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setOrientation(LinearLayout.VERTICAL);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setGravity(Gravity.CENTER);\n", viewVar))
		applyStyle(b, viewVar, v.Style, -1, -1, 0)
		applyCSS(b, viewVar, v.CSS)
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		for _, child := range v.Children {
			buildView(b, child, viewVar, counter)
		}
		return viewVar
	case ir.RowWidget:
		b.WriteString(fmt.Sprintf("        LinearLayout %s = new LinearLayout(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setOrientation(LinearLayout.HORIZONTAL);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setGravity(Gravity.CENTER);\n", viewVar))
		applyStyle(b, viewVar, v.Style, -1, -2, 0)
		applyCSS(b, viewVar, v.CSS)
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		for _, child := range v.Children {
			buildView(b, child, viewVar, counter)
		}
		return viewVar
	case ir.TextViewWidget:
		b.WriteString(fmt.Sprintf("        TextView %s = new TextView(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setText(android.text.Html.fromHtml(\"%s\", android.text.Html.FROM_HTML_MODE_COMPACT));\n", viewVar, v.Text))
		b.WriteString(fmt.Sprintf("        %s.setGravity(Gravity.END);\n", viewVar))
		applyStyle(b, viewVar, v.Style, -1, -2, 0)
		applyTextStyle(b, viewVar, v.Style, 32)
		applyCSS(b, viewVar, v.CSS)
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		return viewVar
	case ir.ButtonWidget:
		b.WriteString(fmt.Sprintf("        Button %s = new Button(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setText(android.text.Html.fromHtml(\"%s\", android.text.Html.FROM_HTML_MODE_COMPACT));\n", viewVar, v.Text))
		applyStyle(b, viewVar, v.Style, 0, -2, 1.0)
		applyTextStyle(b, viewVar, v.Style, 24)
		applyCSS(b, viewVar, v.CSS)
		if v.OnClickFunc != "" {
			// Trigger JNI event
			b.WriteString(fmt.Sprintf("        %s.setOnClickListener(v -> {\n", viewVar))
			b.WriteString(fmt.Sprintf("            NativeBridge.sendEvent(\"%s_onclick\");\n", v.ID))
			b.WriteString(fmt.Sprintf("        });\n"))
		}
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		return viewVar
	case ir.TextFieldWidget:
		b.WriteString(fmt.Sprintf("        EditText %s = new EditText(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setHint(\"%s\");\n", viewVar, v.Placeholder))
		applyStyle(b, viewVar, v.Style, -1, -2, 0)
		applyTextStyle(b, viewVar, v.Style, 24)
		applyCSS(b, viewVar, v.CSS)
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		return viewVar
	case ir.ImageWidget:
		b.WriteString(fmt.Sprintf("        ImageView %s = new ImageView(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        loadImage(%s, \"%s\");\n", viewVar, v.Src))
		applyStyle(b, viewVar, v.Style, -1, -2, 0)
		applyCSS(b, viewVar, v.CSS)
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		return viewVar
	case ir.AudioWidget:
		b.WriteString(fmt.Sprintf("        // AudioWidget has no UI, it plays sound\n"))
		b.WriteString(fmt.Sprintf("        MediaPlayer %s = new MediaPlayer();\n", viewVar))
		b.WriteString("        try {\n")
		b.WriteString(fmt.Sprintf("            %s.setDataSource(\"%s\");\n", viewVar, v.Src))
		b.WriteString(fmt.Sprintf("            %s.prepareAsync();\n", viewVar))
		if v.AutoPlay {
			b.WriteString(fmt.Sprintf("            %s.setOnPreparedListener(mp -> mp.start());\n", viewVar))
		}
		b.WriteString("        } catch (Exception e) { e.printStackTrace(); }\n")
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		// no addView for Audio as it is headless
		return viewVar
	case ir.VideoWidget:
		b.WriteString(fmt.Sprintf("        VideoView %s = new VideoView(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setVideoURI(Uri.parse(\"%s\"));\n", viewVar, v.Src))
		b.WriteString(fmt.Sprintf("        %s.start();\n", viewVar))
		applyStyle(b, viewVar, v.Style, -1, -2, 0)
		applyCSS(b, viewVar, v.CSS)
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		return viewVar
	case ir.WebViewWidget:
		b.WriteString(fmt.Sprintf("        android.webkit.WebView %s = new android.webkit.WebView(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.getSettings().setJavaScriptEnabled(true);\n", viewVar))
		if v.Src != "" {
			b.WriteString(fmt.Sprintf("        %s.loadUrl(\"%s\");\n", viewVar, v.Src))
		} else if v.HTML != "" {
			escapedHTML := strings.ReplaceAll(v.HTML, "\"", "\\\"")
			escapedHTML = strings.ReplaceAll(escapedHTML, "\n", "\\n")
			b.WriteString(fmt.Sprintf("        %s.loadDataWithBaseURL(null, \"%s\", \"text/html\", \"UTF-8\", null);\n", viewVar, escapedHTML))
		}
		applyStyle(b, viewVar, v.Style, -1, -1, 0)
		applyCSS(b, viewVar, v.CSS)
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		return viewVar
	case ir.XMLView:
		b.WriteString(fmt.Sprintf("        int layoutResId_%s = getResources().getIdentifier(\"%s\", \"layout\", getPackageName());\n", viewVar, v.Layout))
		b.WriteString(fmt.Sprintf("        View %s = getLayoutInflater().inflate(layoutResId_%s, null);\n", viewVar, viewVar))
		applyStyle(b, viewVar, v.Style, -1, -2, 0)
		applyCSS(b, viewVar, v.CSS)
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		return viewVar
	case ir.ScrollViewWidget:
		b.WriteString(fmt.Sprintf("        ScrollView %s = new ScrollView(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setFillViewport(true);\n", viewVar))
		applyStyle(b, viewVar, v.Style, -1, -1, 0)
		applyCSS(b, viewVar, v.CSS)
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		for _, child := range v.Children {
			buildView(b, child, viewVar, counter)
		}
		return viewVar
	case ir.CardViewWidget:
		b.WriteString(fmt.Sprintf("        androidx.cardview.widget.CardView %s = new androidx.cardview.widget.CardView(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setRadius(8f);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setCardElevation(4f);\n", viewVar))
		applyStyle(b, viewVar, v.Style, -1, -2, 0)
		applyCSS(b, viewVar, v.CSS)
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		for _, child := range v.Children {
			buildView(b, child, viewVar, counter)
		}
		return viewVar
	case ir.ProgressBarWidget:
		b.WriteString(fmt.Sprintf("        ProgressBar %s = new ProgressBar(this);\n", viewVar))
		applyStyle(b, viewVar, v.Style, -2, -2, 0)
		applyCSS(b, viewVar, v.CSS)
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		return viewVar
	case ir.SwitchWidget:
		b.WriteString(fmt.Sprintf("        Switch %s = new Switch(this);\n", viewVar))
		if v.Checked {
			b.WriteString(fmt.Sprintf("        %s.setChecked(true);\n", viewVar))
		}
		applyStyle(b, viewVar, v.Style, -2, -2, 0)
		applyCSS(b, viewVar, v.CSS)
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		return viewVar
	}

	// Default fallback
	b.WriteString(fmt.Sprintf("        View %s = new View(this);\n", viewVar))
	if parentVar != "" {
		b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
	}
	return viewVar
}
