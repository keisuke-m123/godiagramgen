package renderer

import (
	"unicode"

	"github.com/keisuke-m123/godiagramgen/gocode"
	"github.com/keisuke-m123/godiagramgen/plantuml"
)

type (
	methodRenderer struct{}
)

func newMethodRenderer() *methodRenderer {
	return &methodRenderer{}
}

func (mr *methodRenderer) buildMethods(functions []*gocode.Function) *plantuml.ElementStore {
	elements := plantuml.NewElementStore()
	for _, method := range functions {
		accessModifier := plantuml.AccessModifierPublic
		if unicode.IsLower(rune(method.Name().String()[0])) {
			accessModifier = plantuml.AccessModifierPrivate
		}
		params := make(plantuml.Params, 0)
		for _, p := range method.Parameters() {
			params = append(params, plantuml.Param{Name: p.Name(), Type: p.Type().TypeName().String()})
		}

		returnValues := make(plantuml.ReturnValues, 0)
		for _, r := range method.ReturnValues() {
			returnValues = append(returnValues, plantuml.ReturnValue{
				Type: r.Type().TypeName().String(),
			})
		}

		elements.Add(plantuml.Method(
			accessModifier,
			method.Name().String(),
			params,
			returnValues,
		))
	}
	return elements
}
