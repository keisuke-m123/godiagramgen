package renderer

import (
	"sort"
	"strings"
	"unicode"

	"github.com/keisuke-m123/godiagramgen/gocode"
	"github.com/keisuke-m123/godiagramgen/plantuml"
)

type (
	structRenderer struct {
		relations      *gocode.Relations
		methodRenderer *methodRenderer
	}
)

func newStructRenderer(relations *gocode.Relations) *structRenderer {
	return &structRenderer{
		relations:      relations,
		methodRenderer: newMethodRenderer(),
	}
}

func (r *structRenderer) sortedStructNames(pkgName gocode.PackageName) []gocode.StructName {
	s := r.relations.Structs().PackageStructs(pkgName)
	var names []gocode.StructName
	for i := range s {
		names = append(names, s[i].Name())
	}
	sort.Slice(names, func(i, j int) bool {
		return strings.Compare(names[i].String(), names[j].String()) < 0
	})
	return names
}

func (r *structRenderer) buildStructRelations(pkgName gocode.PackageName) *plantuml.ElementStore {
	elements := plantuml.NewElementStore()
	for _, name := range r.sortedStructNames(pkgName) {
		st, _ := r.relations.Structs().Get(pkgName, name)
		elements.Add(r.buildStructRelation(st).AsSlice()...)
	}
	return elements
}

func (r *structRenderer) buildInPkg(pkgName gocode.PackageName) *plantuml.ElementStore {
	elements := plantuml.NewElementStore()
	for _, name := range r.sortedStructNames(pkgName) {
		s, _ := r.relations.Structs().Get(pkgName, name)
		elements.Add(r.buildElementStructure(s))
	}
	return elements
}

func (r *structRenderer) buildElementStructure(st *gocode.Struct) plantuml.Element {
	fields := r.buildStructFields(st)
	methods := r.buildStructMethods(st)
	color, _ := plantuml.ParseHexColor("#7FFFD4")
	return plantuml.ClassWithOption(
		st.Name().String(),
		plantuml.ClassOptions{
			Spot: plantuml.Spot{
				Name:  'S',
				Color: color,
			},
		},
		fields.Merge(methods).AsSlice()...,
	)
}

func (r *structRenderer) buildStructRelation(st *gocode.Struct) *plantuml.ElementStore {
	elements := plantuml.NewElementStore()
	elements.Add(r.buildCompositions(st).AsSlice()...)
	elements.Add(r.buildExtends(st).AsSlice()...)
	elements.Add(r.buildAggregations(st).AsSlice()...)
	return elements
}

func (r *structRenderer) buildCompositions(s *gocode.Struct) *plantuml.ElementStore {
	elements := plantuml.NewElementStore()
	for _, f := range s.Fields() {
		if f.Embedded() {
			elements.Add(r.buildStructCompositionFromField(s, f).AsSlice()...)
		}
	}
	return elements
}

func (r *structRenderer) buildStructCompositionFromField(structure *gocode.Struct, f *gocode.Field) *plantuml.ElementStore {
	pkgName := structure.PackageSummary().Name().String()
	structName := structure.Name().String()
	elements := plantuml.NewElementStore()

	uniqueTypeNameSet := make(map[string]struct{})
	for _, fType := range append(f.Type().FundamentalTypes(), f.Type()) {
		if _, ok := uniqueTypeNameSet[fType.RelativeFullTypeName().String()]; ok {
			continue
		}
		uniqueTypeNameSet[fType.RelativeFullTypeName().String()] = struct{}{}

		elements.Add(plantuml.Relation(
			plantuml.NewRelationTargetWithNamespace(pkgName, structName),
			plantuml.NewRelationTargetWithNamespace(fType.PackageSummary().Name().String(), fType.TypeName().String()),
			plantuml.RelationTypeComposition,
		))
	}
	return elements
}

func (r *structRenderer) buildAggregations(st *gocode.Struct) *plantuml.ElementStore {
	var orderedFundamentalTypes []*gocode.Type
	for _, f := range st.Fields() {
		if !f.Embedded() {
			orderedFundamentalTypes = append(orderedFundamentalTypes, f.Type().FundamentalTypes()...)
		}
	}
	sort.Slice(orderedFundamentalTypes, func(i, j int) bool {
		return strings.Compare(orderedFundamentalTypes[i].TypeName().String(), orderedFundamentalTypes[j].TypeName().String()) < 0
	})

	elements := plantuml.NewElementStore()
	pkgName := st.PackageSummary().Name().String()
	structName := st.Name().String()

	for _, fType := range orderedFundamentalTypes {
		if fType.Builtin() {
			continue
		}
		fTypePkgName := fType.PackageSummary().Name().String()
		fTypeName := removePointerFromName(fType.TypeName().String())
		elements.Add(plantuml.Relation(
			plantuml.NewRelationTargetWithNamespace(fTypePkgName, fTypeName),
			plantuml.NewRelationTargetWithNamespace(pkgName, structName),
			plantuml.RelationTypeAggregation,
		))
	}

	return elements
}

func (r *structRenderer) buildExtends(st *gocode.Struct) *plantuml.ElementStore {
	elements := plantuml.NewElementStore()
	pkgName := st.PackageSummary().Name().String()
	structName := st.Name().String()
	for _, name := range st.ImplementInterfaces().PackageInterfaceNames() {
		elements.Add(plantuml.Relation(
			plantuml.NewRelationTargetWithNamespace(pkgName, structName),
			plantuml.NewRelationTarget(name.String()),
			plantuml.RelationTypeExtension,
		))
	}
	return elements
}

func (r *structRenderer) buildStructMethods(st *gocode.Struct) *plantuml.ElementStore {
	var orderedFunctions []*gocode.Function
	methods := st.Methods()
	for i := range methods {
		orderedFunctions = append(orderedFunctions, methods[i])
	}
	sort.Slice(orderedFunctions, func(i, j int) bool {
		return strings.Compare(orderedFunctions[i].Name().String(), orderedFunctions[j].Name().String()) < 0
	})

	return r.methodRenderer.buildMethods(orderedFunctions)
}

func (r *structRenderer) buildStructFields(st *gocode.Struct) *plantuml.ElementStore {
	elements := plantuml.NewElementStore()
	for _, field := range st.Fields() {
		accessModifier := plantuml.AccessModifierPublic
		if unicode.IsLower(rune(field.Name().String()[0])) {
			accessModifier = plantuml.AccessModifierPrivate
		}

		elements.Add(plantuml.Field(
			accessModifier,
			field.Name().String(),
			field.Type().RelativeFullTypeName().String(),
		))
	}
	return elements
}
