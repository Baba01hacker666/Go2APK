package ui

const (
	MatchParent = -1
	WrapContent = -2
)

// Style defines styling properties for a widget.
type Style struct {
	BackgroundColor string // e.g. "#FF0000"
	TextColor       string // e.g. "#FFFFFF"
	TextSize        int    // sp
	Padding         int    // dp
	Margin          int    // dp
	Width           int    // MatchParent, WrapContent, or dp
	Height          int    // MatchParent, WrapContent, or dp
	Weight          float32
}
