package pkg

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/keisuke-m123/godiagramgen/testutil"
)

func TestClassDiagram_Render(t *testing.T) {
	tests := []struct {
		name              string
		directories       []string
		ignoreDirectories []string
		theme             string
		wantFilePath      string
	}{
		{
			name:              "GoDiagramGenAll",
			directories:       []string{projectRootPath()},
			ignoreDirectories: []string{path.Join(projectRootPath(), "testingsupport")},
			theme:             "reddress-darkorange",
			wantFilePath:      "../../package-diagram.puml",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, err := NewDiagram(test.directories, test.ignoreDirectories, test.theme)
			if err != nil {
				t.Fatalf("failed newDiagram: %s", err)
			}

			fileBytes, err := ioutil.ReadFile(test.wantFilePath)
			if err != nil {
				t.Fatalf("failed open want file %s: %s", test.wantFilePath, err)
			}

			got := d.Render()
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

func projectRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../../")
}
