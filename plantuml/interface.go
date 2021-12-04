package plantuml

import "fmt"

type (
	iface struct {
		name     string
		elements []Element
	}
)

func (i *iface) Write(builder *LineStringBuilder, indent int) {
	builder.WriteLineWithDepth(indent, fmt.Sprintf("interface %s {", i.name))
	for ei := range i.elements {
		i.elements[ei].Write(builder, indent+1)
	}
	builder.WriteLineWithDepth(indent, "}")
}

func Interface(name string, elements ...Element) Element {
	return &iface{
		name:     name,
		elements: elements,
	}
}
