package ui

// XMLView allows rendering standard Android XML layouts by specifying the layout resource name.
type XMLView struct {
	ID     string
	Layout string // The name of the layout resource (e.g. "my_layout" for res/layout/my_layout.xml)
	Style  Style
	CSS    string
}
