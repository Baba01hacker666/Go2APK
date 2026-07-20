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
        if (currentActivity instanceof MainActivity) {
            currentActivity.runOnUiThread(() ->
                ((MainActivity) currentActivity).updateWidgetText(id, text));
        }
    }

    /** Called from Go to read the current text of a widget. */
    public static String getText(String id) {
        if (currentActivity instanceof MainActivity) {
            return ((MainActivity) currentActivity).getWidgetText(id);
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

func RenderDynamicMainActivity(cfg config.Config, prog *ir.Program) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(`package %s;

import android.app.Activity;
import android.graphics.Color;
import android.os.Bundle;
import android.view.Gravity;
import android.view.View;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.EditText;

public class MainActivity extends Activity {
`, cfg.Package))

	// Declare fields for IDs
	if prog != nil && prog.UI != nil {
		declareFields(&builder, prog.UI)
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

	if prog != nil && prog.UI != nil {
		counter := 0
		rootView := buildView(&builder, prog.UI, "", &counter)
		builder.WriteString(fmt.Sprintf("        return %s;\n", rootView))
	} else {
		// Fallback blank view
		builder.WriteString("        LinearLayout rootView = new LinearLayout(this);\n")
		builder.WriteString("        return rootView;\n")
	}

	builder.WriteString(`    }

    public void updateWidgetText(String id, String text) {
`)
	if prog != nil && prog.UI != nil {
		writeUpdateWidgetCases(&builder, prog.UI)
	}
	builder.WriteString(`    }

    public String getWidgetText(String id) {
`)
	if prog != nil && prog.UI != nil {
		writeGetWidgetCases(&builder, prog.UI)
	}
	builder.WriteString(`        return "";
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
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) this.%s.setText(text); return; }\n", v.ID, v.ID, v.ID))
		}
	case ir.ButtonWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) this.%s.setText(text); return; }\n", v.ID, v.ID, v.ID))
		}
	case ir.TextFieldWidget:
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        if (id.equals(\"%s\")) { if (this.%s != null) this.%s.setText(text); return; }\n", v.ID, v.ID, v.ID))
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
			b.WriteString(fmt.Sprintf("        %s.setBackgroundColor(Color.parseColor(\"%s\"));\n", viewVar, val))
		case "color":
			b.WriteString(fmt.Sprintf("        if (%s instanceof TextView) ((TextView)%s).setTextColor(Color.parseColor(\"%s\"));\n", viewVar, viewVar, val))
		case "font-size":
			size := strings.TrimSuffix(val, "px")
			size = strings.TrimSuffix(size, "dp")
			size = strings.TrimSuffix(size, "sp")
			b.WriteString(fmt.Sprintf("        if (%s instanceof TextView) ((TextView)%s).setTextSize(Float.parseFloat(\"%s\"));\n", viewVar, viewVar, size))
		case "padding":
			pad := strings.TrimSuffix(val, "px")
			b.WriteString(fmt.Sprintf("        %s.setPadding((int)Float.parseFloat(\"%s\"), (int)Float.parseFloat(\"%s\"), (int)Float.parseFloat(\"%s\"), (int)Float.parseFloat(\"%s\"));\n", viewVar, pad, pad, pad, pad))
		case "border-radius":
			rad := strings.TrimSuffix(val, "px")
			b.WriteString(fmt.Sprintf("        {\n"))
			b.WriteString(fmt.Sprintf("            android.graphics.drawable.GradientDrawable gd = new android.graphics.drawable.GradientDrawable();\n"))
			b.WriteString(fmt.Sprintf("            gd.setCornerRadius(Float.parseFloat(\"%s\"));\n", rad))
			b.WriteString(fmt.Sprintf("            %s.setBackground(gd);\n", viewVar))
			b.WriteString(fmt.Sprintf("        }\n"))
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
		b.WriteString(fmt.Sprintf("        %s.setText(\"%s\");\n", viewVar, v.Text))
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
		b.WriteString(fmt.Sprintf("        %s.setText(\"%s\");\n", viewVar, v.Text))
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
	}

	// Default fallback
	b.WriteString(fmt.Sprintf("        View %s = new View(this);\n", viewVar))
	if parentVar != "" {
		b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
	}
	return viewVar
}
