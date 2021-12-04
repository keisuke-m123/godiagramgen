package renderer

import (
	"regexp"
	"sort"
	"strings"

	"github.com/keisuke-m123/godiagramgen/gocode"
	"github.com/keisuke-m123/godiagramgen/plantuml"
)

type RenderingOptions struct {
	Title string
	Notes string
	Theme string
}

type Renderer struct {
	relations           *gocode.Relations
	renderingOptions    *RenderingOptions
	structRenderer      *structRenderer
	interfaceRenderer   *interfaceRenderer
	definedTypeRenderer *definedTypeRenderer
	aliasRenderer       *aliasRenderer
}

func NewRenderer(relations *gocode.Relations, options *RenderingOptions) *Renderer {
	return &Renderer{
		relations:           relations,
		renderingOptions:    options,
		structRenderer:      newStructRenderer(relations),
		interfaceRenderer:   newInterfaceRenderer(relations),
		definedTypeRenderer: newDefinedTypeRenderer(relations),
		aliasRenderer:       newAliasRender(relations),
	}
}

func (r *Renderer) Render() *plantuml.Result {
	elements := plantuml.NewElementStore()
	if r.renderingOptions.Theme != "" {
		elements.Add(plantuml.Theme(r.renderingOptions.Theme))
	}
	if r.renderingOptions.Title != "" {
		elements.Add(plantuml.Title(r.renderingOptions.Title))
	}
	if note := strings.TrimSpace(r.renderingOptions.Notes); note != "" {
		elements.Add(plantuml.Legend(r.renderingOptions.Notes))
	}

	var packageNames []gocode.PackageName
	for _, pkg := range r.relations.Packages().AsSlice() {
		packageNames = append(packageNames, pkg.Summary().Name())
	}
	sort.Slice(packageNames, func(i, j int) bool {
		return strings.Compare(packageNames[i].String(), packageNames[j].String()) < 0
	})
	for _, pkgName := range packageNames {
		elements.Add(r.buildPackage(pkgName).AsSlice()...)
	}
	elements.Add(r.aliasRenderer.buildTypeAliases().AsSlice()...)
	return plantuml.PlantUML(elements.AsSlice()...)
}

func (r *Renderer) buildPackage(pkgName gocode.PackageName) *plantuml.ElementStore {
	structs := r.structRenderer.buildInPkg(pkgName)
	interfaces := r.interfaceRenderer.buildInPkg(pkgName)
	definedTypes := r.definedTypeRenderer.buildInPkg(pkgName)

	elements := plantuml.NewElementStore()
	elements.Add(plantuml.Namespace(
		pkgName.String(),
		structs.Merge(interfaces).Merge(definedTypes).AsSlice()...,
	))

	elements.Add(r.structRenderer.buildStructRelations(pkgName).AsSlice()...)
	elements.Add(r.interfaceRenderer.buildRelations(pkgName).AsSlice()...)
	elements.Add(r.definedTypeRenderer.buildRelations(pkgName).AsSlice()...)

	return elements
}

func generateRenamedName(currentName string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(currentName, "")
}

func removePointerFromName(name string) string {
	reg, _ := regexp.Compile(`\*+`)
	return reg.ReplaceAllString(name, "")
}
