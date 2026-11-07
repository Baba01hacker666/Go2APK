package android

import "testing"

func TestJavaStringLiteral(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain", "Hello", `"Hello"`},
		{"double quote", `press "me"`, `"press \"me\""`},
		{"backslash", `a\b`, `"a\\b"`},
		{"newline", "a\nb", `"a\nb"`},
		{"tab and cr", "a\tb\rc", `"a\tb\rc"`},
		{"unicode passes through", "héllo ✓", `"héllo ✓"`},
		{"control char uses octal, not \\u", "a\x07b", `"a\007b"`},
		{"empty", "", `""`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := javaStringLiteral(tc.in); got != tc.want {
				t.Errorf("javaStringLiteral(%q) = %s, want %s", tc.in, got, tc.want)
			}
		})
	}
}

// Java expands \uXXXX escapes before lexing, so emitting one inside a string
// literal would insert a real newline and break the generated source. Control
// characters must therefore never come out as \u escapes.
func TestJavaStringLiteralAvoidsUnicodeEscapes(t *testing.T) {
	for r := rune(0); r < 0x20; r++ {
		got := javaStringLiteral(string(r))
		if contains(got, "\\u") {
			t.Errorf("javaStringLiteral(%q) = %s, which contains a \\u escape", r, got)
		}
	}
	got := javaStringLiteral("\x7f")
	if contains(got, "\\u") {
		t.Errorf("javaStringLiteral(DEL) = %s, which contains a \\u escape", got)
	}
}

func TestXMLAttrEscaping(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{`android.permission.CAMERA`, `android.permission.CAMERA`},
		{`a & b`, `a &amp; b`},
		{`<script>`, `&lt;script&gt;`},
		{`say "hi"`, `say &quot;hi&quot;`},
		{"line\nbreak", "line&#10;break"},
	}
	for _, tc := range tests {
		if got := xmlAttr(tc.in); got != tc.want {
			t.Errorf("xmlAttr(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
