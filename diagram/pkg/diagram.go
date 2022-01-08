package pkg

import (
	"fmt"

	"github.com/keisuke-m123/goanalyzer/gocode"
	"github.com/spf13/afero"
)

type (
	Diagram struct {
		renderer *renderer
	}
)

func NewDiagram(directoryPaths []string, ignoreDirectories []string, theme string) (*Diagram, error) {
	loadOptions := &gocode.LoadOptions{
		Directories:        directoryPaths,
		IgnoredDirectories: ignoreDirectories,
		Recursive:          true,
		FileSystem:         afero.NewOsFs(),
	}
	relations, err := gocode.LoadRelations(loadOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to load relations: %w", err)
	}

	return &Diagram{
		renderer: newRenderer(theme, relations.GeneratePackageGraph()),
	}, nil
}

func (d *Diagram) Render() string {
	return d.renderer.render()
}
