package ui

const (
	MatchParent = -1
	WrapContent = -2
)

// Animation defines basic view animation properties.
type Animation struct {
	Type     string // e.g., "fade_in", "slide_up", "bounce", "pulse"
	Duration int    // duration in milliseconds
	Delay    int    // delay in milliseconds
	Loop     bool   // whether the animation repeats
}

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
	Animation       Animation // Optional animation to play when view is shown
}
