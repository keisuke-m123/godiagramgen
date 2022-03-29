package classdiagram

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	goplantuml "github.com/keisuke-m123/godiagramgen/diagram/class"
	"github.com/keisuke-m123/godiagramgen/diagram/class/renderer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	FlagIgnore                 = "ignore"
	FlagTitle                  = "title"
	FlagNotes                  = "notes"
	FlagOutput                 = "output"
	FlagTheme                  = "theme"
	FlagRecursive              = "recursive"
	FlagRenderExternalPackages = "render-external-packages"
)

type FlagValues struct {
	Ignore                 string
	Title                  string
	Notes                  string
	Output                 string
	Theme                  string
	Recursive              bool
	RenderExternalPackages bool
}

type FlagSet struct {
	set    *pflag.FlagSet
	values FlagValues
}

func (fs *FlagSet) InitializeFlags() {
	s := fs.set
	vs := &fs.values
	s.StringVar(&vs.Ignore, FlagIgnore, "", "Comma separated list of folders to ignore")
	s.StringVar(&vs.Title, FlagTitle, "", "Title of the generated diagram")
	s.StringVar(&vs.Notes, FlagNotes, "", "Comma separated list of notes to be added to the diagram")
	s.StringVar(&vs.Output, FlagOutput, "", "Output file path. If omitted, then this will default to standard output")
	s.StringVar(&vs.Theme, FlagTheme, "", "Change theme")
	s.BoolVar(&vs.Recursive, FlagRecursive, false, "Walk all directories recursively")
	s.BoolVar(&vs.RenderExternalPackages, FlagRenderExternalPackages, false, "Render external packages")
}

func (fs *FlagSet) Values() FlagValues {
	return fs.values
}

func NewClassDiagramGenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "class",
		Short: "generate class diagram from specified packages",
	}

	fs := FlagSet{set: cmd.PersistentFlags()}
	fs.InitializeFlags()

	cmd.Run = func(cmd *cobra.Command, args []string) { run(fs.Values(), args) }

	return cmd
}

func run(flagValues FlagValues, args []string) {
	var noteList []string
	if flagValues.Notes != "" {
		noteList = append(noteList, "<b><u>Notes</u></b>")
	}
	split := strings.Split(flagValues.Notes, ",")
	for _, note := range split {
		trimmed := strings.TrimSpace(note)
		if trimmed != "" {
			noteList = append(noteList, trimmed)
		}
	}
	renderingOptions := &renderer.RenderingOptions{
		Title:                  flagValues.Title,
		Notes:                  strings.Join(noteList, "\n"),
		Theme:                  flagValues.Theme,
		RenderExternalPackages: flagValues.RenderExternalPackages,
	}

	dirs, err := getDirectories(args)
	if err != nil {
		fmt.Println("usage:\ngoplantuml <DIR>\nDIR Must be a valid directory")
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	ignoredDirectories, err := getIgnoredDirectories(flagValues.Ignore)
	if err != nil {
		fmt.Println("usage:\ngoplantuml [-ignore=<DIRLIST>]\nDIRLIST Must be a valid comma separated list of existing directories")
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cd, err := goplantuml.NewDiagram(dirs, ignoredDirectories, flagValues.Recursive, renderingOptions)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	rendered := cd.Render()
	var writer io.Writer
	if flagValues.Output != "" {
		writer, err = os.Create(flagValues.Output)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
		}
	} else {
		writer = os.Stdout
	}
	_, _ = fmt.Fprint(writer, rendered)
}

func getDirectories(args []string) ([]string, error) {
	if len(args) < 1 {
		return nil, errors.New("DIR missing")
	}
	var dirs []string
	for _, dir := range args {
		fi, err := os.Stat(dir)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("could not find directory %s", dir)
		}
		if !fi.Mode().IsDir() {
			return nil, fmt.Errorf("%s is not a directory", dir)
		}
		dirAbs, err := filepath.Abs(dir)
		if err != nil {
			return nil, fmt.Errorf("could not find directory %s", dir)
		}
		dirs = append(dirs, dirAbs)
	}
	return dirs, nil
}

func getIgnoredDirectories(list string) ([]string, error) {
	var result []string
	list = strings.TrimSpace(list)
	if list == "" {
		return result, nil
	}
	split := strings.Split(list, ",")
	for _, dir := range split {
		dirAbs, err := filepath.Abs(strings.TrimSpace(dir))
		if err != nil {
			return nil, fmt.Errorf("could not find directory %s", dir)
		}
		result = append(result, dirAbs)
	}
	return result, nil
}
