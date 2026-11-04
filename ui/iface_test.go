package ui

// Compile-time check that every exported widget satisfies Widget.
var (
	_ Widget = Column{}
	_ Widget = Row{}
	_ Widget = TextView{}
	_ Widget = Button{}
	_ Widget = TextField{}
	_ Widget = Image{}
	_ Widget = Audio{}
	_ Widget = Video{}
	_ Widget = WebView{}
	_ Widget = ScrollView{}
	_ Widget = CardView{}
	_ Widget = ProgressBar{}
	_ Widget = Switch{}
	_ Widget = XMLView{}
)
