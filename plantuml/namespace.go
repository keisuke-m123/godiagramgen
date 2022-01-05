package plantuml

import "fmt"

type (
	namespace struct {
		val      string
		as       string
		elements []Element
	}

	NamespaceOptions struct {
		As string
	}
)

func (n *namespace) Write(builder *LineStringBuilder, indent int) {
	var as string
	if n.as != "" {
		as = fmt.Sprintf(" as %s", n.as)
	}

	if as == "" {
		builder.WriteLineWithDepth(indent, fmt.Sprintf(`namespace %s {`, n.val))
	} else {
		builder.WriteLineWithDepth(indent, fmt.Sprintf(`namespace %s %s {`, n.val, as))
	}

	for i := range n.elements {
		n.elements[i].Write(builder, indent+1)
	}
	builder.WriteLineWithDepth(indent, "}")
}

func Namespace(val string, elements ...Element) Element {
	return &namespace{
		val:      val,
		elements: elements,
	}
}

func NamespaceWithOption(val string, options NamespaceOptions, elements ...Element) Element {
	return &namespace{
		val:      val,
		elements: elements,
		as:       options.As,
	}
}
