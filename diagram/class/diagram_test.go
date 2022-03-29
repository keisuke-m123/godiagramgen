package class

import (
	"io/ioutil"
	"testing"

	"github.com/keisuke-m123/godiagramgen/diagram/class/renderer"
	"github.com/keisuke-m123/godiagramgen/testutil"
)

func TestClassDiagram_Render(t *testing.T) {
	tests := []struct {
		name              string
		renderingOptions  *renderer.RenderingOptions
		recursive         bool
		directories       []string
		ignoreDirectories []string
		wantFilePath      string
	}{
		{
			name:             "TestingSupportAll",
			renderingOptions: &renderer.RenderingOptions{Theme: "reddress-darkorange"},
			recursive:        true,
			directories:      []string{"../../testingsupport"},
			wantFilePath:     "../../testingsupport/testingsupport-all.puml",
		},
		{
			name:             "TestingSupportAllIgnoreDirectories",
			renderingOptions: &renderer.RenderingOptions{},
			recursive:        true,
			directories:      []string{"../../testingsupport"},
			ignoreDirectories: []string{
				"../../testingsupport/subfolder",
				"../../testingsupport/subfolder2",
				"../../testingsupport/connectionlabels",
			},
			wantFilePath: "../../testingsupport/testingsupport-all-ignore-directories.puml",
		},
		{
			name: "TestingSupport",
			renderingOptions: &renderer.RenderingOptions{
				Title: "Test Title",
				Notes: "<b><u>Notes</u></b>\nExample 1\nExample 1 continues\nExample 2",
			},
			directories:  []string{"../../testingsupport"},
			wantFilePath: "../../testingsupport/testingsupport.puml",
		},
		{
			name: "TestingSupportWithRenderExternalPackages",
			renderingOptions: &renderer.RenderingOptions{
				RenderExternalPackages: true,
			},
			directories:  []string{"../../testingsupport"},
			wantFilePath: "../../testingsupport/testingsupport-render-external-packages.puml",
		},
		{
			name:             "ParenTheSizedTypeDeclarations",
			renderingOptions: &renderer.RenderingOptions{},
			directories:      []string{"../../testingsupport", "../../testingsupport/parenthesizedtypedeclarations"},
			wantFilePath:     "../../testingsupport/testingsupport-parenthesizedtypedeclarations.puml",
		},
		{
			name:             "AliasMethods",
			renderingOptions: &renderer.RenderingOptions{},
			directories:      []string{"../../testingsupport/aliasmethods"},
			wantFilePath:     "../../testingsupport/aliasmethods.puml",
		},
		{
			name:             "SubFolder1-3",
			renderingOptions: &renderer.RenderingOptions{},
			directories: []string{
				"../../testingsupport/subfolder",
				"../../testingsupport/subfolder2",
				"../../testingsupport/subfolder3",
			},
			wantFilePath: "../../testingsupport/subfolder1-3.puml",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, err := NewDiagram(
				test.directories,
				test.ignoreDirectories,
				test.recursive,
				test.renderingOptions,
			)
			if err != nil {
				t.Fatalf("failed newDiagramWithOptions: %s", err)
			}

			renderResult := d.Render()
			fileBytes, err := ioutil.ReadFile(test.wantFilePath)
			if err != nil {
				t.Fatalf("failed open want file %s: %s", test.wantFilePath, err)
			}

			got := renderResult.String()
			want := string(fileBytes)

			if got != want {
				t.Errorf(
					"failed render: want %s\n\ngot %s\n\ndiff: %s",
					want,
					got,
					testutil.Diff(t, want, got),
				)
			}
		})
	}
}
