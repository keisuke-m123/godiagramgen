package renderer

import (
	"fmt"
	"sort"

	"github.com/keisuke-m123/goanalyzer/gocode"
	"github.com/keisuke-m123/godiagramgen/plantuml"
)

type (
	aliasRenderer struct {
		relations *gocode.Relations
	}
)

func newAliasRender(relations *gocode.Relations) *aliasRenderer {
	return &aliasRenderer{
		relations: relations,
	}
}

func (ar *aliasRenderer) buildTypeAliases() *plantuml.ElementStore {
	orderedAliases := ar.relations.TypeAliases().AliasAll()
	sort.Slice(orderedAliases, func(i, j int) bool {
		ai := orderedAliases[i]
		aj := orderedAliases[j]
		return fmt.Sprintf("%s %s %s", ai.Name(), ai.PackageSummary().Name(), ai.Type().TypeName()) <
			fmt.Sprintf("%s %s %s", aj.Name(), aj.PackageSummary().Name(), aj.Type().TypeName())
	})

	elements := plantuml.NewElementStore()
	color, _ := plantuml.ParseHexColor("#EDDC44")
	for _, alias := range orderedAliases {
		pkgName := alias.PackageSummary().Name().String()
		aliasName := alias.Name().String()
		var stereotype string
		if alias.Type().Builtin() {
			stereotype = fmt.Sprintf("alias of __%s__", alias.Type().TypeName().String())
		}
		elements.Add(
			plantuml.Namespace(
				pkgName,
				plantuml.ClassWithOption(
					aliasName,
					plantuml.ClassOptions{
						Stereotype: plantuml.Stereotype(stereotype),
						Spot: plantuml.Spot{
							Name:  'T',
							Color: color,
						},
					},
				),
			),
		)
		if alias.Type().Builtin() {
			continue
		}

		rtFrom := plantuml.NewRelationTargetWithNamespace(pkgName, aliasName)
		rtTo := plantuml.NewRelationTargetWithNamespace(
			alias.Type().PackageSummary().Name().String(),
			alias.Type().TypeName().String(),
		)

		elements.Add(plantuml.Relation(rtFrom, rtTo, plantuml.RelationTypeAlias))
	}

	return elements
}
