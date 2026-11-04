package android

import (
	"fmt"
	"strings"
)

// javaStringLiteral returns s formatted as a Java string literal, including the
// surrounding double quotes.
//
// Every character that would otherwise terminate the literal or form an invalid
// escape sequence is escaped, so arbitrary Go string values — embedded quotes,
// backslashes, newlines and non-ASCII text — always produce compilable Java.
func javaStringLiteral(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 2)
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		default:
			// Java expands \uXXXX escapes before lexing, so a \u000a escape
			// would be turned into a real newline and break the literal.
			// Octal escapes are handled by the lexer and are therefore safe.
			if r < 0x20 || r == 0x7f {
				fmt.Fprintf(&b, `\%03o`, r)
				continue
			}
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}

// xmlAttr escapes s for use inside a double-quoted XML attribute value.
func xmlAttr(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&apos;")
		case '\n':
			b.WriteString("&#10;")
		case '\r':
			b.WriteString("&#13;")
		case '\t':
			b.WriteString("&#9;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// xmlText escapes s for use as XML character data.
func xmlText(s string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	).Replace(s)
}
