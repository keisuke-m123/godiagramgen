package plantuml

import "fmt"

type theme struct {
	val string
}

func (t *theme) Write(builder *LineStringBuilder, indent int) {
	builder.WriteLineWithDepth(indent, fmt.Sprintf("!theme %s", t.val))
	t.themeAdjustment(builder, indent)
}

func (t *theme) themeAdjustment(builder *LineStringBuilder, indent int) {
	switch t.val {
	case "reddress-darkblue",
		"reddress-darkgreen",
		"reddress-darkorange",
		"reddress-darkred",
		"reddress-lightblue",
		"reddress-lightgreen",
		"reddress-lightorange",
		"reddress-lightred":
		builder.WriteLineWithDepth(indent, "skinparam class {")
		builder.WriteLineWithDepth(indent+1, "attributeIconSize 8")
		builder.WriteLineWithDepth(indent, "}")
	}
}

func Theme(val string) Element {
	return &theme{val: val}
}
