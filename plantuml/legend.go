package plantuml

type legend struct {
	note string
}

func (l *legend) Write(builder *LineStringBuilder, indent int) {
	builder.WriteLineWithDepth(indent, "legend")
	builder.WriteLineWithDepth(indent, l.note)
	builder.WriteLineWithDepth(indent, "end legend")
}

func Legend(note string) Element {
	return &legend{note: note}
}
