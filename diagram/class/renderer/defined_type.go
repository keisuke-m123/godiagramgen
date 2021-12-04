package renderer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/keisuke-m123/godiagramgen/gocode"
	"github.com/keisuke-m123/godiagramgen/plantuml"
)

type (
	definedTypeRenderer struct {
		relations      *gocode.Relations
		methodRenderer *methodRenderer
	}
)

func newDefinedTypeRenderer(relations *gocode.Relations) *definedTypeRenderer {
	return &definedTypeRenderer{
		relations:      relations,
		methodRenderer: newMethodRenderer(),
	}
}

func (r *definedTypeRenderer) buildRelations(pkgName gocode.PackageName) *plantuml.ElementStore {
	elements := plantuml.NewElementStore()
	for _, name := range r.sortedDefinedTypeNames(pkgName) {
		dt, _ := r.relations.DefinedTypes().Get(pkgName, name)
		elements.Add(r.buildRelation(pkgName, dt).AsSlice()...)
	}
	return elements
}

func (r *definedTypeRenderer) buildRelation(pkgName gocode.PackageName, dt *gocode.DefinedType) *plantuml.ElementStore {
	elements := plantuml.NewElementStore()

	typ := dt.Type()
	if typ.Builtin() {
		return elements
	}

	typeName := typ.TypeName().String()
	if renamed := generateRenamedName(typeName); typeName != renamed {
		if typ.ContainsBuiltinInFundamentalTypes() {
			typeName = renamed
		}
		elements.Add(newTypedClass(pkgName.String(), typ.TypeName().String(), typeName))
	}

	rtFrom := plantuml.NewRelationTargetWithNamespace(pkgName.String(), dt.Name().String())
	rtTo := plantuml.NewRelationTargetWithNamespace(typ.PackageSummary().Name().String(), typeName)
	elements.Add(plantuml.Relation(rtFrom, rtTo, plantuml.RelationTypeAlias))

	return elements
}

func (r *definedTypeRenderer) buildInPkg(pkgName gocode.PackageName) *plantuml.ElementStore {
	elements := plantuml.NewElementStore()
	for _, name := range r.sortedDefinedTypeNames(pkgName) {
		definedType, _ := r.relations.DefinedTypes().Get(pkgName, name)
		elements.Add(r.build(definedType))
	}
	return elements
}

func (r *definedTypeRenderer) build(definedType *gocode.DefinedType) plantuml.Element {
	color, _ := plantuml.ParseHexColor("#FF7700")
	var stereotype string
	if definedType.Type().Builtin() {
		stereotype = fmt.Sprintf("type of __%s__", definedType.Type().TypeName().String())
	}

	return plantuml.ClassWithOption(
		definedType.Name().String(),
		plantuml.ClassOptions{
			Stereotype: plantuml.Stereotype(stereotype),
			Spot: plantuml.Spot{
				Name:  'D',
				Color: color,
			},
		},
		r.buildMethods(definedType).AsSlice()...,
	)
}

func (r *definedTypeRenderer) buildMethods(definedType *gocode.DefinedType) *plantuml.ElementStore {
	var orderedFunctions []*gocode.Function
	methods := definedType.Methods()
	for i := range methods {
		orderedFunctions = append(orderedFunctions, methods[i])
	}
	sort.Slice(orderedFunctions, func(i, j int) bool {
		return strings.Compare(orderedFunctions[i].Name().String(), orderedFunctions[j].Name().String()) < 0
	})
	return r.methodRenderer.buildMethods(orderedFunctions)
}

func (r *definedTypeRenderer) sortedDefinedTypeNames(pkgName gocode.PackageName) []gocode.DefinedTypeName {
	definedTypes := r.relations.DefinedTypes().PackageDefinedTypes(pkgName)
	var names []gocode.DefinedTypeName
	for i := range definedTypes {
		names = append(names, definedTypes[i].Name())
	}
	sort.Slice(names, func(i, j int) bool {
		return strings.Compare(names[i].String(), names[j].String()) < 0
	})
	return names
}
