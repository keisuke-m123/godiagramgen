package pkg

import (
	"strings"

	"github.com/keisuke-m123/goanalyzer/gocode"
	"github.com/keisuke-m123/godiagramgen/plantuml"
)

type renderer struct {
	theme    string
	pkgGraph *gocode.PackageGraph
}

func newRenderer(theme string, pkgGraph *gocode.PackageGraph) *renderer {
	return &renderer{
		theme:    theme,
		pkgGraph: pkgGraph,
	}
}

func (r *renderer) render() string {
	elements := plantuml.NewElementStore()
	if r.theme != "" {
		elements.Add(plantuml.Theme(r.theme))
	}
	for _, path := range r.pkgGraph.SortedPackagePaths() {
		elements.Add(r.buildNamespace(path))
		for _, imPath := range r.pkgGraph.SortedImportPackagePaths(path) {
			elements.Add(r.buildNamespace(imPath.Path()))
		}
	}
	for _, path := range r.pkgGraph.SortedPackagePaths() {
		for _, imPath := range r.pkgGraph.SortedImportPackagePaths(path) {
			elements.Add(plantuml.Relation(
				plantuml.NewRelationTarget(r.relationTargetName(path)),
				plantuml.NewRelationTarget(r.relationTargetName(imPath.Path())),
				plantuml.RelationTypeArrow,
			))
		}
	}
	return plantuml.PlantUML(elements.AsSlice()...).String()
}

func (r *renderer) buildNamespace(pkgPath gocode.PackagePath) plantuml.Element {
	ps := strings.Split(r.namespacePath(pkgPath), "/")
	var ns plantuml.Element
	for i := len(ps) - 1; i >= 0; i-- {
		if ns == nil {
			ns = plantuml.Namespace(ps[i])
		} else {
			ns = plantuml.Namespace(ps[i], ns)
		}
	}
	return ns
}

func (r *renderer) namespacePath(pkgPath gocode.PackagePath) string {
	return strings.ReplaceAll(pkgPath.String(), ".", "")
}

func (r *renderer) relationTargetName(pkgPath gocode.PackagePath) string {
	t := r.namespacePath(pkgPath)
	return strings.ReplaceAll(t, "/", ".")
}
