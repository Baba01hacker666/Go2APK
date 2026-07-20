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
    }

    private View createAppView() {
`, cfg.Name))

	if prog != nil && prog.UI != nil {
		buildView(&builder, prog.UI, "rootView")
		builder.WriteString("        return rootView;\n")
	} else {
		// Fallback blank view
		builder.WriteString("        LinearLayout rootView = new LinearLayout(this);\n")
		builder.WriteString("        return rootView;\n")
	}

	builder.WriteString(`    }
}
`)
	return builder.String()
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

var viewCounter int

func buildView(b *strings.Builder, w ir.Widget, parentVar string) string {
	viewCounter++
	viewVar := fmt.Sprintf("view%d", viewCounter)

	switch v := w.(type) {
	case ir.ColumnWidget:
		b.WriteString(fmt.Sprintf("        LinearLayout %s = new LinearLayout(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setOrientation(LinearLayout.VERTICAL);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setGravity(Gravity.CENTER);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setPadding(24, 24, 24, 24);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setLayoutParams(new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.MATCH_PARENT));\n", viewVar))
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "rootView" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		} else {
			// Actually the first one is the rootView, so we reassign
			b.WriteString(fmt.Sprintf("        LinearLayout %s = %s;\n", parentVar, viewVar))
		}
		for _, child := range v.Children {
			buildView(b, child, viewVar)
		}
		return viewVar
	case ir.RowWidget:
		b.WriteString(fmt.Sprintf("        LinearLayout %s = new LinearLayout(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setOrientation(LinearLayout.HORIZONTAL);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setGravity(Gravity.CENTER);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setLayoutParams(new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT));\n", viewVar))
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" && parentVar != "rootView" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		for _, child := range v.Children {
			buildView(b, child, viewVar)
		}
		return viewVar
	case ir.TextViewWidget:
		b.WriteString(fmt.Sprintf("        TextView %s = new TextView(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setText(\"%s\");\n", viewVar, v.Text))
		b.WriteString(fmt.Sprintf("        %s.setTextSize(32);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setGravity(Gravity.END);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setPadding(16, 16, 16, 16);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setLayoutParams(new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT));\n", viewVar))
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" && parentVar != "rootView" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		return viewVar
	case ir.ButtonWidget:
		b.WriteString(fmt.Sprintf("        Button %s = new Button(this);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setText(\"%s\");\n", viewVar, v.Text))
		b.WriteString(fmt.Sprintf("        %s.setTextSize(24);\n", viewVar))
		b.WriteString(fmt.Sprintf("        LinearLayout.LayoutParams lp_%s = new LinearLayout.LayoutParams(0, LinearLayout.LayoutParams.WRAP_CONTENT, 1.0f);\n", viewVar))
		b.WriteString(fmt.Sprintf("        lp_%s.setMargins(8, 8, 8, 8);\n", viewVar))
		b.WriteString(fmt.Sprintf("        %s.setLayoutParams(lp_%s);\n", viewVar, viewVar))
		if v.OnClickFunc != "" {
			// Trigger JNI event
			b.WriteString(fmt.Sprintf("        %s.setOnClickListener(v -> {\n", viewVar))
			b.WriteString(fmt.Sprintf("            NativeBridge.sendEvent(\"%s_onclick\");\n", v.ID))
			b.WriteString(fmt.Sprintf("        });\n"))
		}
		if v.ID != "" {
			b.WriteString(fmt.Sprintf("        this.%s = %s;\n", v.ID, viewVar))
		}
		if parentVar != "" && parentVar != "rootView" {
			b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
		}
		return viewVar
	}

	// Default fallback
	b.WriteString(fmt.Sprintf("        View %s = new View(this);\n", viewVar))
	if parentVar != "" && parentVar != "rootView" {
		b.WriteString(fmt.Sprintf("        %s.addView(%s);\n", parentVar, viewVar))
	}
	return viewVar
}
