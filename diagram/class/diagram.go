package class

import (
	"fmt"

	"github.com/keisuke-m123/godiagramgen/diagram/class/renderer"
	"github.com/keisuke-m123/godiagramgen/gocode"
	"github.com/keisuke-m123/godiagramgen/plantuml"
	"github.com/spf13/afero"
)

type Diagram struct {
	renderer *renderer.Renderer
}

func newDiagramWithOptions(
	loadOptions *gocode.LoadOptions,
	renderingOptions *renderer.RenderingOptions,
) (*Diagram, error) {
	r, err := gocode.LoadRelations(loadOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to load relations: %w", err)
	}
	return &Diagram{renderer: renderer.NewRenderer(r, renderingOptions)}, nil
}

func NewDiagram(
	directoryPaths []string,
	ignoreDirectories []string,
	recursive bool,
	renderingOptions *renderer.RenderingOptions,
) (*Diagram, error) {
	loadOptions := &gocode.LoadOptions{
		Directories:        directoryPaths,
		IgnoredDirectories: ignoreDirectories,
		Recursive:          recursive,
		FileSystem:         afero.NewOsFs(),
	}
	return newDiagramWithOptions(loadOptions, renderingOptions)
}

func (d *Diagram) Render() *plantuml.Result {
	return d.renderer.Render()
}
