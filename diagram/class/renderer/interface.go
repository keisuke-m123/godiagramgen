package renderer

import (
	"sort"
	"strings"

	"github.com/keisuke-m123/goanalyzer/gocode"
	"github.com/keisuke-m123/godiagramgen/plantuml"
)

type (
	interfaceRenderer struct {
		relations      *gocode.Relations
		methodRenderer *methodRenderer
	}
)

func newInterfaceRenderer(relations *gocode.Relations) *interfaceRenderer {
	return &interfaceRenderer{
		relations:      relations,
		methodRenderer: newMethodRenderer(),
	}
}

func (r *interfaceRenderer) buildInPkg(pkgName gocode.PackageName) *plantuml.ElementStore {
	elements := plantuml.NewElementStore()
	for _, name := range r.sortedInterfaceNames(pkgName) {
		iface, _ := r.relations.Interfaces().Get(pkgName, name)
		elements.Add(r.buildInterface(iface))
	}
	return elements
}

func (r *interfaceRenderer) buildRelations(pkgName gocode.PackageName) *plantuml.ElementStore {
	elements := plantuml.NewElementStore()
	for _, name := range r.sortedInterfaceNames(pkgName) {
		iface, _ := r.relations.Interfaces().Get(pkgName, name)
		elements.Add(r.buildCompositions(iface).AsSlice()...)
	}
	return elements
}

func (r *interfaceRenderer) buildInterface(iface *gocode.Interface) plantuml.Element {
	methods := r.buildMethods(iface)
	return plantuml.Interface(iface.Name().String(), methods.AsSlice()...)
}

func (r *interfaceRenderer) buildMethods(iface *gocode.Interface) *plantuml.ElementStore {
	var orderedFunctions []*gocode.Function
	methods := iface.Methods()
	for i := range methods {
		orderedFunctions = append(orderedFunctions, methods[i])
	}
	sort.Slice(orderedFunctions, func(i, j int) bool {
		return strings.Compare(orderedFunctions[i].Name().String(), orderedFunctions[j].Name().String()) < 0
	})

	return r.methodRenderer.buildMethods(orderedFunctions)
}

func (r *interfaceRenderer) buildCompositions(iface *gocode.Interface) *plantuml.ElementStore {
	var orderedEmbeds []*gocode.Embed
	embeds := iface.Embeds()
	for i := range embeds {
		orderedEmbeds = append(orderedEmbeds, embeds[i])
	}
	sort.Slice(orderedEmbeds, func(i, j int) bool {
		return strings.Compare(
			orderedEmbeds[i].Type().TypeName().String(),
			orderedEmbeds[j].Type().TypeName().String(),
		) < 0
	})

	elements := plantuml.NewElementStore()
	for i := range orderedEmbeds {
		e := orderedEmbeds[i]
		elements.Add(plantuml.Relation(
			plantuml.NewRelationTargetWithNamespace(iface.PackageSummary().Name().String(), iface.Name().String()),
			plantuml.NewRelationTargetWithNamespace(e.Type().PackageSummary().Name().String(), e.Type().TypeName().String()),
			plantuml.RelationTypeComposition,
		))
	}
	return elements
}

func (r *interfaceRenderer) sortedInterfaceNames(pkgName gocode.PackageName) []gocode.InterfaceName {
	interfaces := r.relations.Interfaces().PackageInterfaces(pkgName)
	var names []gocode.InterfaceName
	for i := range interfaces {
		names = append(names, interfaces[i].Name())
	}
	sort.Slice(names, func(i, j int) bool {
		return strings.Compare(names[i].String(), names[j].String()) < 0
	})
	return names
}
