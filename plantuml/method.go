package plantuml

import (
	"fmt"
	"strings"
)

type (
	method struct {
		accessModifier AccessModifier
		name           string
		parameters     Params
		returnValues   ReturnValues
	}

	Params []Param

	Param struct {
		Name string
		Type string
	}

	ReturnValues []ReturnValue

	ReturnValue struct {
		Name string
		Type string
	}
)

func (p Param) toString() string {
	if p.Name != "" {
		return fmt.Sprintf("%s %s", p.Name, p.Type)
	}
	return p.Type
}

func (ps Params) toString() string {
	var ss []string
	for i := range ps {
		ss = append(ss, ps[i].toString())
	}
	return strings.Join(ss, ", ")
}

func (rv ReturnValue) toString() string {
	if rv.Name != "" {
		return fmt.Sprintf("%s %s", rv.Name, rv.Type)
	}
	return rv.Type
}

func (rvs ReturnValues) toString() string {
	var ss []string
	for i := range rvs {
		ss = append(ss, rvs[i].toString())
	}
	return strings.Join(ss, ", ")
}

func (m *method) buildParameters() string {
	return m.parameters.toString()
}

func (m *method) buildReturnValues() string {
	if len(m.returnValues) > 1 {
		return fmt.Sprintf("(%s)", m.returnValues.toString())
	}
	return m.returnValues.toString()
}

func (m *method) Write(builder *LineStringBuilder, indent int) {
	builder.WriteLineWithDepth(indent, fmt.Sprintf(
		"%s %s(%s) %s",
		m.accessModifier.toString(),
		m.name,
		m.buildParameters(),
		m.buildReturnValues(),
	))
}

func Method(
	accessModifier AccessModifier,
	name string,
	parameters Params,
	returnValues ReturnValues,
) Element {
	return &method{
		accessModifier: accessModifier,
		name:           name,
		parameters:     parameters,
		returnValues:   returnValues,
	}
}
