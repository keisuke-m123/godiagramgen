package plantuml

import (
	"fmt"
	"strings"
)

const (
	AccessModifierPublic AccessModifier = iota
	AccessModifierPrivate
)

type (
	AccessModifier int

	class struct {
		name       string
		elements   []Element
		as         string
		stereotype Stereotype
		spot       Spot
	}

	ClassOptions struct {
		As         string
		Spot       Spot
		Stereotype Stereotype
	}

	Spot struct {
		Name  rune
		Color *Color
	}

	Stereotype string
)

func (am AccessModifier) toString() string {
	switch am {
	case AccessModifierPublic:
		return "+"
	case AccessModifierPrivate:
		return "-"
	default:
		return ""
	}
}

func (s Spot) build() string {
	if s.Name != 0 && s.Color != nil {
		return fmt.Sprintf("(%s,  %s)", string(s.Name), s.Color.HexRGBA())
	}
	return ""
}

func (s Stereotype) build() string {
	return strings.TrimSpace(string(s))
}

func Class(name string, elements ...Element) Element {
	return ClassWithOption(name, ClassOptions{}, elements...)
}

func ClassWithOption(name string, options ClassOptions, elements ...Element) Element {
	return &class{
		name:       name,
		elements:   elements,
		as:         options.As,
		stereotype: options.Stereotype,
		spot:       options.Spot,
	}
}

func (c *class) Write(builder *LineStringBuilder, indent int) {
	var as string
	if c.as != "" {
		as = fmt.Sprintf(`as %s`, c.as)
	}

	builder.WriteLineWithDepth(indent, fmt.Sprintf(`class "%s" %s %s {`, c.name, as, c.buildStereotype()))
	for i := range c.elements {
		c.elements[i].Write(builder, indent+1)
	}
	builder.WriteLineWithDepth(indent, "}")
}

func (c *class) buildStereotype() string {
	sp := c.spot.build()
	st := c.stereotype.build()
	if sp != "" || st != "" {
		return fmt.Sprintf("<< %s %s >>", sp, st)
	}
	return ""
}
