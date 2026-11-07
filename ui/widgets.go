package ui

// XMLView renders a native Android XML layout from res/layout by resource name.
// It lets an app drop down to hand-written layouts when the declarative widget
// set cannot express something.
type XMLView struct {
	ID     string
	Layout string // Layout resource name, e.g. "my_layout" for res/layout/my_layout.xml
	Style  Style
	CSS    string
}

func (XMLView) isWidget() {}
