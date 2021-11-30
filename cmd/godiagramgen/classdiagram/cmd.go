package classdiagram

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	goplantuml "github.com/keisuke-m123/godiagramgen/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	FlagIgnore                  = "ignore"
	FlagTitle                   = "title"
	FlagNotes                   = "notes"
	FlagOutput                  = "output"
	FlagTheme                   = "theme"
	FlagRecursive               = "recursive"
	FlagAggregatePrivateMembers = "aggregate-private-members"
	FlagShowOptionsAsNote       = "show-options-as-notes"
	FlagShowAggregations        = "show-aggregations"
	FlagShowCompositions        = "show-compositions"
	FlagShowImplementations     = "show-implementations"
	FlagShowAliases             = "show-aliases"
	FlagShowConnectionLabels    = "show-connection-labels"
	FlagHideFields              = "hide-fields"
	FlagHideMethods             = "hide-methods"
	FlagHideConnections         = "hide-connections"
	FlagHidePrivateMembers      = "hide-private-members"
)

type FlagValues struct {
	Ignore                  string
	Title                   string
	Notes                   string
	Output                  string
	Theme                   string
	Recursive               bool
	AggregatePrivateMembers bool
	ShowAggregations        bool
	ShowCompositions        bool
	ShowImplementations     bool
	ShowAliases             bool
	ShowConnectionLabels    bool
	ShowOptionsAsNote       bool
	HideFields              bool
	HideMethods             bool
	HideConnections         bool
	HidePrivateMembers      bool
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
	s.BoolVar(&vs.AggregatePrivateMembers, FlagAggregatePrivateMembers, false, "Show aggregations for private members. Ignored if -show-aggregations is not used.")
	s.BoolVar(&vs.ShowAggregations, FlagShowAggregations, false, "Renders public aggregations even when -hide-connections is used (do not render by default)")
	s.BoolVar(&vs.ShowCompositions, FlagShowCompositions, false, "Shows compositions even when -hide-connections is used")
	s.BoolVar(&vs.ShowImplementations, FlagShowImplementations, false, "Shows implementations even when -hide-connections is used")
	s.BoolVar(&vs.ShowAliases, FlagShowAliases, false, "Shows aliases even when -hide-connections is used")
	s.BoolVar(&vs.ShowConnectionLabels, FlagShowConnectionLabels, false, "Shows labels in the connections to identify the connections types (e.g. extends, implements, aggregates, alias of")
	s.BoolVar(&vs.ShowOptionsAsNote, FlagShowOptionsAsNote, false, "Show a note in the diagram with the none evident options ran with this CLI")
	s.BoolVar(&vs.HideFields, FlagHideFields, false, "Hides fields")
	s.BoolVar(&vs.HideMethods, FlagHideMethods, false, "Hides methods")
	s.BoolVar(&vs.HideConnections, FlagHideConnections, false, "Hides all connections in the diagram")
	s.BoolVar(&vs.HidePrivateMembers, FlagHidePrivateMembers, false, "Hide private fields and methods")
}

func (fs *FlagSet) Values() FlagValues {
	return fs.values
}

// RenderingOptionSlice will implements the sort interface
type RenderingOptionSlice []goplantuml.RenderingOption

// Len is the number of elements in the collection.
func (as RenderingOptionSlice) Len() int {
	return len(as)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (as RenderingOptionSlice) Less(i, j int) bool {
	return as[i] < as[j]
}

// Swap swaps the elements with indexes i and j.
func (as RenderingOptionSlice) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
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
	renderingOptions := map[goplantuml.RenderingOption]interface{}{
		goplantuml.RenderConnectionLabels:  flagValues.ShowConnectionLabels,
		goplantuml.RenderFields:            !flagValues.HideFields,
		goplantuml.RenderMethods:           !flagValues.HideMethods,
		goplantuml.RenderAggregations:      flagValues.ShowAggregations,
		goplantuml.RenderTitle:             flagValues.Title,
		goplantuml.AggregatePrivateMembers: flagValues.AggregatePrivateMembers,
		goplantuml.RenderPrivateMembers:    !flagValues.HidePrivateMembers,
		goplantuml.RenderTheme:             flagValues.Theme,
	}
	if flagValues.HideConnections {
		renderingOptions[goplantuml.RenderAliases] = flagValues.ShowAliases
		renderingOptions[goplantuml.RenderCompositions] = flagValues.ShowCompositions
		renderingOptions[goplantuml.RenderImplementations] = flagValues.ShowImplementations
	}

	var noteList []string
	if flagValues.ShowOptionsAsNote {
		legend, err := getLegend(renderingOptions)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		noteList = append(noteList, legend)
	}
	if flagValues.Notes != "" {
		noteList = append(noteList, "", "<b><u>Notes</u></b>")
	}
	split := strings.Split(flagValues.Notes, ",")
	for _, note := range split {
		trimmed := strings.TrimSpace(note)
		if trimmed != "" {
			noteList = append(noteList, trimmed)
		}
	}
	renderingOptions[goplantuml.RenderNotes] = strings.Join(noteList, "\n")

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

	result, err := goplantuml.NewClassDiagram(dirs, ignoredDirectories, flagValues.Recursive)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := result.SetRenderingOptions(renderingOptions); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	rendered := result.Render()
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

func getLegend(ro map[goplantuml.RenderingOption]interface{}) (string, error) {
	result := "<u><b>Legend</b></u>\n"
	orderedOptions := RenderingOptionSlice{}
	for o := range ro {
		orderedOptions = append(orderedOptions, o)
	}
	sort.Sort(orderedOptions)
	for _, option := range orderedOptions {
		val := ro[option]
		switch option {
		case goplantuml.RenderAggregations:
			result = fmt.Sprintf("%sRender Aggregations: %t\n", result, val.(bool))
		case goplantuml.RenderAliases:
			result = fmt.Sprintf("%sRender Connections: %t\n", result, val.(bool))
		case goplantuml.RenderCompositions:
			result = fmt.Sprintf("%sRender Compositions: %t\n", result, val.(bool))
		case goplantuml.RenderFields:
			result = fmt.Sprintf("%sRender Fields: %t\n", result, val.(bool))
		case goplantuml.RenderImplementations:
			result = fmt.Sprintf("%sRender Implementations: %t\n", result, val.(bool))
		case goplantuml.RenderMethods:
			result = fmt.Sprintf("%sRender Methods: %t\n", result, val.(bool))
		case goplantuml.AggregatePrivateMembers:
			result = fmt.Sprintf("%sPritave Aggregations: %t\n", result, val.(bool))
		}
	}
	return strings.TrimSpace(result), nil
}
