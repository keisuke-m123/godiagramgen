package plantuml

import "fmt"

type title struct {
	val string
}

func (t *title) Write(builder *LineStringBuilder, indent int) {
	builder.WriteLineWithDepth(indent, fmt.Sprintf("title %s", t.val))
}

func Title(val string) Element {
	return &title{val: val}
}
