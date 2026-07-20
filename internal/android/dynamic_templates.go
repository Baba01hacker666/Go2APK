package android

import (
	"fmt"
	"strings"

	"github.com/Baba01hacker666/Go2APK/internal/config"
	"github.com/Baba01hacker666/Go2APK/ir"
)

// RenderNativeBridge creates the generic Java JNI bridge.
func RenderNativeBridge(cfg config.Config) string {
	return fmt.Sprintf(`package %s;

import android.app.Activity;

final class NativeBridge {
    private static final String LIBRARY_NAME = "go2apkapp";

    static {
        System.loadLibrary(LIBRARY_NAME);
    }

    private NativeBridge() {
    }

    static void sendEvent(String eventName) {
        sendEventToGo(eventName);
    }

    private static native void sendEventToGo(String eventName);

    public static void setActivity(Activity activity) {
        currentActivity = activity;
    }

    private static Activity currentActivity;

    public static void updateText(String id, String text) {
        if (currentActivity instanceof MainActivity) {
            currentActivity.runOnUiThread(() -> {
                ((MainActivity) currentActivity).updateWidgetText(id, text);
            });
        }
    }
}
`, cfg.Package)
}

// RenderDynamicMainActivity creates a Java Activity by traversing the IR widget tree.
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

func buildView(b *strings.Builder, w ir.Widget, parentVar string, counter *int) string {
	*counter++
	viewVar := fmt.Sprintf("view%d", *counter)

	switch v := w.(type) {
	case ir.ColumnWidget:
		b.WriteString(fmt.Sprintf("        LinearLayout %s = new LinearLayout(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setOrientation(LinearLayout.VERTICAL);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setGravity(Gravity.CENTER);\n", viewVar))
		applyStyle(b, viewVar, v.Style, -1, -1, 0)
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
