package plantuml

import "fmt"

type (
	field struct {
		accessModifier AccessModifier
		name           string
		typ            string
	}
)

func (f *field) Write(builder *LineStringBuilder, indent int) {
	s := fmt.Sprintf("%s %s %s", f.accessModifier.toString(), f.name, f.typ)
	builder.WriteLineWithDepth(indent, s)
}

func Field(accessModifier AccessModifier, name, typ string) Element {
	return &field{
		accessModifier: accessModifier,
		name:           name,
		typ:            typ,
	}
}
