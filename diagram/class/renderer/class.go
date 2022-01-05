package renderer

import (
	"strings"

	"github.com/keisuke-m123/godiagramgen/plantuml"
)

func newTypedClass(namespace, name, as string) plantuml.Element {
	spotName := '-'
	switch {
	case strings.HasPrefix(name, "[]"):
		spotName = 's'
	case strings.HasPrefix(name, "map"):
		spotName = 'm'
	case strings.HasPrefix(name, "func"):
		spotName = 'f'
	}
	color, _ := plantuml.ParseHexColor("#3CB371")
	return plantuml.Namespace(
		namespace,
		plantuml.ClassWithOption(
			name,
			plantuml.ClassOptions{
				As: as,
				Spot: plantuml.Spot{
					Name:  spotName,
					Color: color,
				},
			},
		),
	)
}
