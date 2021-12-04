package plantuml

import (
	"strings"
)

const (
	tab = "    "
)

// LineStringBuilder extends the strings.Builder and adds functionality to build a string with tabs and
// adding new lines
type LineStringBuilder struct {
	strings.Builder
}

func newLineStringBuilder() *LineStringBuilder {
	return &LineStringBuilder{}
}

// WriteLineWithDepth will write the given text with added tabs at the beginning into the string builder.
func (lsb *LineStringBuilder) WriteLineWithDepth(depth int, str string) {
	_, _ = lsb.WriteString(strings.Repeat(tab, depth))
	_, _ = lsb.WriteString(str)
	_, _ = lsb.WriteString("\n")
}

type Result struct {
	builder *LineStringBuilder
}

func (r *Result) String() string {
	return r.builder.String()
}

func newResult(builder *LineStringBuilder) *Result {
	return &Result{builder: builder}
}

type Element interface {
	Write(builder *LineStringBuilder, indent int)
}

type ElementStore struct {
	elements []Element
}

func NewElementStore() *ElementStore {
	return &ElementStore{}
}

func (e *ElementStore) Add(es ...Element) {
	e.elements = append(e.elements, es...)
}

func (e *ElementStore) Merge(es *ElementStore) *ElementStore {
	res := NewElementStore()
	res.Add(append(e.elements, es.elements...)...)
	return res
}

func (e *ElementStore) AsSlice() []Element {
	return e.elements
}

func PlantUML(elements ...Element) *Result {
	builder := newLineStringBuilder()
	builder.WriteLineWithDepth(0, "@startuml")
	for i := range elements {
		elements[i].Write(builder, 0)
	}
	builder.WriteLineWithDepth(0, "@enduml")
	return newResult(builder)
}
