/*
Package parser generates PlantUml http://plantuml.com/ Class diagrams for your golang projects
The main structure is the ClassParser which you can generate by calling the NewClassDiagram(dir)
function.

Pass the directory where the .go files are and the parser will analyze the code and build a structure
containing the information it needs to Render the class diagram.

call the Render() function and this will return a string with the class diagram.
*/
package parser

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/spf13/afero"
	"golang.org/x/tools/go/packages"
)

// LineStringBuilder extends the strings.Builder and adds functionality to build a string with tabs and
// adding new lines
type LineStringBuilder struct {
	strings.Builder
}

const tab = "    "
const builtinPackageName = "__builtin__"
const implements = `"implements"`
const extends = `"extends"`
const aggregates = `"uses"`
const aliasOf = `"alias of"`

// WriteLineWithDepth will write the given text with added tabs at the beginning into the string builder.
func (lsb *LineStringBuilder) WriteLineWithDepth(depth int, str string) {
	_, _ = lsb.WriteString(strings.Repeat(tab, depth))
	_, _ = lsb.WriteString(str)
	_, _ = lsb.WriteString("\n")
}

// ClassDiagramOptions will provide a way for callers of the NewClassDiagramFs() function to pass all the necessary arguments.
type ClassDiagramOptions struct {
	FileSystem         afero.Fs
	Directories        []string
	IgnoredDirectories []string
	RenderingOptions   map[RenderingOption]interface{}
	Recursive          bool
}

// RenderingOptions will allow the class parser to optionally enebale or disable the things to render.
type RenderingOptions struct {
	Title                   string
	Notes                   string
	Theme                   string
	Aggregations            bool
	Fields                  bool
	Methods                 bool
	Compositions            bool
	Implementations         bool
	Aliases                 bool
	ConnectionLabels        bool
	AggregatePrivateMembers bool
	PrivateMembers          bool
}

const aliasComplexNameComment = "'This class was created so that we can correctly have an alias pointing to this name. Since it contains dots that can break namespaces"

const (
	// RenderAggregations is to be used in the SetRenderingOptions argument as the key to the map, when value is true, it will set the parser to render aggregations
	RenderAggregations RenderingOption = iota

	// RenderCompositions is to be used in the SetRenderingOptions argument as the key to the map, when value is true, it will set the parser to render compositions
	RenderCompositions

	// RenderImplementations is to be used in the SetRenderingOptions argument as the key to the map, when value is true, it will set the parser to render implementations
	RenderImplementations

	// RenderAliases is to be used in the SetRenderingOptions argument as the key to the map, when value is true, it will set the parser to render aliases
	RenderAliases

	// RenderFields is to be used in the SetRenderingOptions argument as the key to the map, when value is true, it will set the parser to render fields
	RenderFields

	// RenderMethods is to be used in the SetRenderingOptions argument as the key to the map, when value is true, it will set the parser to render methods
	RenderMethods

	// RenderConnectionLabels is to be used in the SetRenderingOptions argument as the key to the map, when value is true, it will set the parser to render the connection labels
	RenderConnectionLabels

	// RenderTitle is the options for the Title of the diagram. The value of this will be rendered as a title unless empty
	RenderTitle

	// RenderNotes contains a list of notes to be rendered in the class diagram
	RenderNotes

	// AggregatePrivateMembers is to be used in the SetRenderingOptions argument as the key to the map, when value is true, it will connect aggregations with private members
	AggregatePrivateMembers

	// RenderPrivateMembers is used if private members (fields, methods) should be rendered
	RenderPrivateMembers

	// RenderTheme is the options for the Theme of the diagram. The value of this will be rendered as a "none" theme if empty
	RenderTheme
)

// RenderingOption is an alias for an it so it is easier to use it as options in a map (see SetRenderingOptions(map[RenderingOption]bool) error)
type RenderingOption int

// ClassParser contains the structure of the parsed files. The structure is a map of package_names that contains
// a map of structure_names -> Structs
type ClassParser struct {
	renderingOptions  *RenderingOptions
	structure         map[string]map[string]*Struct
	allInterfaces     map[string]struct{}
	allStructs        map[string]struct{}
	allImports        map[string]string
	allAliases        map[string]*Alias
	allRenamedStructs map[string]map[string]string
}

// NewClassDiagramWithOptions returns a new classParser with which can Render the class diagram of
// files in the given directory passed in the ClassDiargamOptions. This will also alow for different types of FileSystems
// Passed since it is part of the ClassDiagramOptions as well.
func NewClassDiagramWithOptions(options *ClassDiagramOptions) (*ClassParser, error) {
	classParser := &ClassParser{
		renderingOptions: &RenderingOptions{
			Aggregations:     false,
			Fields:           true,
			Methods:          true,
			Compositions:     true,
			Implementations:  true,
			Aliases:          true,
			ConnectionLabels: false,
			Title:            "",
			Notes:            "",
		},
		structure:         make(map[string]map[string]*Struct),
		allInterfaces:     make(map[string]struct{}),
		allStructs:        make(map[string]struct{}),
		allImports:        make(map[string]string),
		allAliases:        make(map[string]*Alias),
		allRenamedStructs: make(map[string]map[string]string),
	}
	ignoreDirectoryMap := map[string]struct{}{}
	for _, dir := range options.IgnoredDirectories {
		ignoreDirectoryMap[dir] = struct{}{}
	}
	for _, directoryPath := range options.Directories {
		if options.Recursive {
			err := afero.Walk(options.FileSystem, directoryPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					if strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor" {
						return filepath.SkipDir
					}
					if _, ok := ignoreDirectoryMap[path]; ok {
						return filepath.SkipDir
					}
					_ = classParser.parseDirectory(path)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		} else {
			err := classParser.parseDirectory(directoryPath)
			if err != nil {
				return nil, err
			}
		}
	}

	for s := range classParser.allStructs {
		st := classParser.getStruct(s)
		if st != nil {
			for i := range classParser.allInterfaces {
				inter := classParser.getStruct(i)
				if st.ImplementsInterface(inter) {
					st.AddToExtends(i)
				}
			}
		}
	}
	_ = classParser.SetRenderingOptions(options.RenderingOptions)
	return classParser, nil
}

// NewClassDiagram returns a new classParser with which can Render the class diagram of
// files in the given directory
func NewClassDiagram(directoryPaths []string, ignoreDirectories []string, recursive bool) (*ClassParser, error) {
	options := &ClassDiagramOptions{
		Directories:        directoryPaths,
		IgnoredDirectories: ignoreDirectories,
		Recursive:          recursive,
		RenderingOptions:   map[RenderingOption]interface{}{},
		FileSystem:         afero.NewOsFs(),
	}
	return NewClassDiagramWithOptions(options)
}

// parse the given *packages.Package into the ClassParser structure
func (p *ClassParser) parsePackage(pkg *packages.Package) {
	_, ok := p.structure[pkg.Name]
	if !ok {
		p.structure[pkg.Name] = make(map[string]*Struct)
	}
	p.parseTypes(pkg)
	p.parseFunctions(pkg)
}

func (p *ClassParser) parseImports(pkg *packages.Package) {
	// general imports
	for packageName, importedPkg := range pkg.Imports {
		p.allImports[importedPkg.Name] = packageName
	}
	// aliases
	for _, d := range pkg.TypesInfo.Defs {
		pkgName, ok := d.(*types.PkgName)
		if ok {
			p.allImports[pkgName.Name()] = pkgName.Imported().Path()
		}
	}
}

func (p *ClassParser) parseTypes(pkg *packages.Package) {
	for _, name := range pkg.Types.Scope().Names() {
		obj := pkg.Types.Scope().Lookup(name)
		switch {
		case p.parseStruct(pkg.Name, obj):
		case p.parseInterface(pkg.Name, obj):
		case p.parseAlias(pkg.Name, obj):
		}
	}
}

func (p *ClassParser) parseFunctions(pkg *packages.Package) {
	for i := range pkg.TypesInfo.Defs {
		d := pkg.TypesInfo.Defs[i]
		funcObj, ok := d.(*types.Func)
		if !ok {
			continue
		}
		tFunc, ok := funcObj.Type().Underlying().(*types.Signature)
		if !ok {
			continue
		}
		if tFunc.Recv() == nil || tFunc.Recv().Name() == "" {
			continue
		}

		recvTypeName := getFieldType(tFunc.Recv().Type(), p.allImports)
		recvTypeName = replacePackageConstant(recvTypeName, "")
		if recvTypeName[0] == "*"[0] {
			recvTypeName = recvTypeName[1:]
		}
		structure := p.getOrCreateStruct(pkg.Name, recvTypeName)
		if structure.Type == "" {
			structure.Type = "class"
		}
		fullName := fmt.Sprintf("%s.%s", pkg.Name, recvTypeName)
		p.allStructs[fullName] = struct{}{}
		structure.AddMethod(funcObj, p.allImports)
	}
}

func (p *ClassParser) parseStruct(currentPackageName string, obj types.Object) (parsed bool) {
	tStruct, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return false
	}
	for i := 0; i < tStruct.NumFields(); i++ {
		p.getOrCreateStruct(currentPackageName, obj.Name()).AddField(tStruct.Field(i), p.allImports)
	}
	p.getOrCreateStruct(currentPackageName, obj.Name()).Type = "class"
	p.allStructs[fmt.Sprintf("%s.%s", currentPackageName, obj.Name())] = struct{}{}
	return true
}

func (p *ClassParser) parseInterface(currentPackageName string, obj types.Object) (parsed bool) {
	tInterface, ok := obj.Type().Underlying().(*types.Interface)
	if !ok {
		return false
	}
	for i := 0; i < tInterface.NumMethods(); i++ {
		m := tInterface.Method(i)
		p.getOrCreateStruct(currentPackageName, obj.Name()).AddMethod(m, p.allImports)
	}
	for i := 0; i < tInterface.NumEmbeddeds(); i++ {
		e := tInterface.EmbeddedType(i)
		tName := getFieldType(e, p.allImports)
		p.getOrCreateStruct(currentPackageName, obj.Name()).AddToComposition(replacePackageConstant(tName, currentPackageName))
	}
	p.getOrCreateStruct(currentPackageName, obj.Name()).Type = "interface"
	p.allInterfaces[fmt.Sprintf("%s.%s", currentPackageName, obj.Name())] = struct{}{}
	return true
}

func (p *ClassParser) parseAlias(currentPackageName string, obj types.Object) (parsed bool) {
	if _, ok := obj.(*types.TypeName); !ok {
		return false
	}
	typeName := obj.Name()
	if !isPrimitiveString(typeName) {
		typeName = fmt.Sprintf("%s.%s", currentPackageName, typeName)
	}
	p.getOrCreateStruct(currentPackageName, typeName).Type = "alias"

	packageName := currentPackageName
	aliasType := getFieldType(obj.Type().Underlying(), p.allImports)
	aliasType = replacePackageConstant(aliasType, "")
	if isPrimitiveString(aliasType) {
		packageName = builtinPackageName
	}

	alias := getNewAlias(fmt.Sprintf("%s.%s", packageName, aliasType), currentPackageName, typeName)
	p.allAliases[typeName] = alias
	if strings.Count(alias.Name, ".") > 1 {
		pack := strings.SplitN(alias.Name, ".", 2)
		if _, ok := p.allRenamedStructs[pack[0]]; !ok {
			p.allRenamedStructs[pack[0]] = map[string]string{}
		}
		renamedClass := generateRenamedStructName(pack[1])
		p.allRenamedStructs[pack[0]][renamedClass] = pack[1]
	}

	return true
}

func (p *ClassParser) parseDirectory(directoryPath string) error {
	loadConfig := &packages.Config{
		Mode: packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedSyntax |
			packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports,
		Dir: directoryPath,
	}
	pkgs, err := packages.Load(loadConfig)
	if err != nil {
		return fmt.Errorf("load packages failed: %w", err)
	}
	for i := range pkgs {
		p.parseImports(pkgs[i])
		p.parsePackage(pkgs[i])
	}
	return nil
}

// Render returns a string of the class diagram that this parser has generated.
func (p *ClassParser) Render() string {
	str := &LineStringBuilder{}
	str.WriteLineWithDepth(0, "@startuml")
	if p.renderingOptions.Theme != "" {
		str.WriteLineWithDepth(0, fmt.Sprintf("!theme %s", p.renderingOptions.Theme))
		str.WriteLineWithDepth(0, themeAdjustment(p.renderingOptions.Theme))
	}
	if p.renderingOptions.Title != "" {
		str.WriteLineWithDepth(0, fmt.Sprintf(`title %s`, p.renderingOptions.Title))
	}
	if note := strings.TrimSpace(p.renderingOptions.Notes); note != "" {
		str.WriteLineWithDepth(0, "legend")
		str.WriteLineWithDepth(0, note)
		str.WriteLineWithDepth(0, "end legend")
	}

	var packs []string
	for pack := range p.structure {
		packs = append(packs, pack)
	}
	sort.Strings(packs)
	for _, pack := range packs {
		structures := p.structure[pack]
		p.renderStructures(pack, structures, str)
	}
	if p.renderingOptions.Aliases {
		p.renderAliases(str)
	}
	if !p.renderingOptions.Fields {
		str.WriteLineWithDepth(0, "hide fields")
	}
	if !p.renderingOptions.Methods {
		str.WriteLineWithDepth(0, "hide methods")
	}
	str.WriteLineWithDepth(0, "@enduml")
	return str.String()
}

func (p *ClassParser) renderStructures(pack string, structures map[string]*Struct, str *LineStringBuilder) {
	if len(structures) > 0 {
		composition := &LineStringBuilder{}
		extends := &LineStringBuilder{}
		aggregations := &LineStringBuilder{}
		str.WriteLineWithDepth(0, fmt.Sprintf(`namespace %s {`, pack))

		var names []string
		for name := range structures {
			names = append(names, name)
		}

		sort.Strings(names)

		for _, name := range names {
			structure := structures[name]
			p.renderStructure(structure, name, str, composition, extends, aggregations)
		}
		var orderedRenamedStructs []string
		for tempName := range p.allRenamedStructs[pack] {
			orderedRenamedStructs = append(orderedRenamedStructs, tempName)
		}
		sort.Strings(orderedRenamedStructs)
		for _, tempName := range orderedRenamedStructs {
			name := p.allRenamedStructs[pack][tempName]
			str.WriteLineWithDepth(1, fmt.Sprintf(`class "%s" as %s {`, name, tempName))
			str.WriteLineWithDepth(2, aliasComplexNameComment)
			str.WriteLineWithDepth(1, "}")
		}
		str.WriteLineWithDepth(0, fmt.Sprintf(`}`))
		if p.renderingOptions.Compositions {
			str.WriteLineWithDepth(0, composition.String())
		}
		if p.renderingOptions.Implementations {
			str.WriteLineWithDepth(0, extends.String())
		}
		if p.renderingOptions.Aggregations {
			str.WriteLineWithDepth(0, aggregations.String())
		}
	}
}

func (p *ClassParser) renderAliases(str *LineStringBuilder) {
	aliasString := ""
	if p.renderingOptions.ConnectionLabels {
		aliasString = aliasOf
	}
	orderedAliases := AliasSlice{}
	for _, alias := range p.allAliases {
		orderedAliases = append(orderedAliases, *alias)
	}
	sort.Sort(orderedAliases)
	for _, alias := range orderedAliases {
		aliasName := alias.Name
		if strings.Count(alias.Name, ".") > 1 {
			split := strings.SplitN(alias.Name, ".", 2)
			if aliasRename, ok := p.allRenamedStructs[split[0]]; ok {
				renamed := generateRenamedStructName(split[1])
				if _, ok := aliasRename[renamed]; ok {
					aliasName = fmt.Sprintf("%s.%s", split[0], renamed)
				}
			}
		}
		str.WriteLineWithDepth(0, fmt.Sprintf(`"%s" #.. %s"%s"`, aliasName, aliasString, alias.AliasOf))
	}
}

func (p *ClassParser) renderStructure(structure *Struct, name string, str *LineStringBuilder, composition *LineStringBuilder, extends *LineStringBuilder, aggregations *LineStringBuilder) {
	privateFields := &LineStringBuilder{}
	publicFields := &LineStringBuilder{}
	privateMethods := &LineStringBuilder{}
	publicMethods := &LineStringBuilder{}
	sType := ""
	renderStructureType := structure.Type
	switch structure.Type {
	case "class":
		sType = "<< (S,Aquamarine) >>"
	case "alias":
		sType = "<< (T, #FF7700) >> "
		renderStructureType = "class"
	}
	str.WriteLineWithDepth(1, fmt.Sprintf(`%s %s %s {`, renderStructureType, name, sType))
	p.renderStructFields(structure, privateFields, publicFields)
	p.renderStructMethods(structure, privateMethods, publicMethods)
	p.renderCompositions(structure, name, composition)
	p.renderExtends(structure, name, extends)
	p.renderAggregations(structure, name, aggregations)
	if privateFields.Len() > 0 {
		str.WriteLineWithDepth(0, privateFields.String())
	}
	if publicFields.Len() > 0 {
		str.WriteLineWithDepth(0, publicFields.String())
	}
	if privateMethods.Len() > 0 {
		str.WriteLineWithDepth(0, privateMethods.String())
	}
	if publicMethods.Len() > 0 {
		str.WriteLineWithDepth(0, publicMethods.String())
	}
	str.WriteLineWithDepth(1, fmt.Sprintf(`}`))
}

func (p *ClassParser) renderCompositions(structure *Struct, name string, composition *LineStringBuilder) {
	var orderedCompositions []string

	for c := range structure.Composition {
		if !strings.Contains(c, ".") {
			c = fmt.Sprintf("%s.%s", p.getPackageName(c, structure), c)
		}
		composedString := ""
		if p.renderingOptions.ConnectionLabels {
			composedString = extends
		}
		c = fmt.Sprintf(`"%s" *-- %s"%s.%s"`, c, composedString, structure.PackageName, name)
		orderedCompositions = append(orderedCompositions, c)
	}
	sort.Strings(orderedCompositions)
	for _, c := range orderedCompositions {
		composition.WriteLineWithDepth(0, c)
	}
}

func (p *ClassParser) renderAggregations(structure *Struct, name string, aggregations *LineStringBuilder) {

	aggregationMap := structure.Aggregations
	if p.renderingOptions.AggregatePrivateMembers {
		p.updatePrivateAggregations(structure, aggregationMap)
	}
	p.renderAggregationMap(aggregationMap, structure, aggregations, name)
}

func (p *ClassParser) updatePrivateAggregations(structure *Struct, aggregationsMap map[string]struct{}) {
	for agg := range structure.PrivateAggregations {
		aggregationsMap[agg] = struct{}{}
	}
}

func (p *ClassParser) renderAggregationMap(aggregationMap map[string]struct{}, structure *Struct, aggregations *LineStringBuilder, name string) {
	var orderedAggregations []string
	for a := range aggregationMap {
		orderedAggregations = append(orderedAggregations, a)
	}

	sort.Strings(orderedAggregations)

	for _, a := range orderedAggregations {
		if !strings.Contains(a, ".") {
			a = fmt.Sprintf("%s.%s", p.getPackageName(a, structure), a)
		}
		aggregationString := ""
		if p.renderingOptions.ConnectionLabels {
			aggregationString = aggregates
		}
		if p.getPackageName(a, structure) != builtinPackageName {
			aggregations.WriteLineWithDepth(0, fmt.Sprintf(`"%s.%s"%s o-- "%s"`, structure.PackageName, name, aggregationString, a))
		}
	}
}

func (p *ClassParser) getPackageName(t string, st *Struct) string {
	packageName := st.PackageName
	if isPrimitiveString(t) {
		packageName = builtinPackageName
	}
	return packageName
}
func (p *ClassParser) renderExtends(structure *Struct, name string, extends *LineStringBuilder) {
	var orderedExtends []string
	for c := range structure.Extends {
		if !strings.Contains(c, ".") {
			c = fmt.Sprintf("%s.%s", structure.PackageName, c)
		}
		implementString := ""
		if p.renderingOptions.ConnectionLabels {
			implementString = implements
		}
		c = fmt.Sprintf(`"%s" <|-- %s"%s.%s"`, c, implementString, structure.PackageName, name)
		orderedExtends = append(orderedExtends, c)
	}
	sort.Strings(orderedExtends)
	for _, c := range orderedExtends {
		extends.WriteLineWithDepth(0, c)
	}
}

func (p *ClassParser) renderStructMethods(structure *Struct, privateMethods *LineStringBuilder, publicMethods *LineStringBuilder) {
	var orderedFunctions []*Function
	for i := range structure.Functions {
		orderedFunctions = append(orderedFunctions, structure.Functions[i])
	}
	sort.Slice(orderedFunctions, func(i, j int) bool {
		return strings.Compare(orderedFunctions[i].Name, orderedFunctions[j].Name) < 0
	})

	for _, method := range orderedFunctions {
		accessModifier := "+"
		if unicode.IsLower(rune(method.Name[0])) {
			if !p.renderingOptions.PrivateMembers {
				continue
			}

			accessModifier = "-"
		}
		parameterList := make([]string, 0)
		for _, p := range method.Parameters {
			parameterList = append(parameterList, fmt.Sprintf("%s %s", p.Name, p.Type))
		}
		returnValues := ""
		if len(method.ReturnValues) > 0 {
			if len(method.ReturnValues) == 1 {
				returnValues = fmt.Sprintf(" %s", method.ReturnValues[0])
			} else {
				returnValues = fmt.Sprintf(" (%s)", strings.Join(method.ReturnValues, ", "))
			}
		}
		if accessModifier == "-" {
			privateMethods.WriteLineWithDepth(2, fmt.Sprintf(`%s %s(%s)%s`, accessModifier, method.Name, strings.Join(parameterList, ", "), returnValues))
		} else {
			publicMethods.WriteLineWithDepth(2, fmt.Sprintf(`%s %s(%s)%s`, accessModifier, method.Name, strings.Join(parameterList, ", "), returnValues))
		}
	}
}

func (p *ClassParser) renderStructFields(structure *Struct, privateFields *LineStringBuilder, publicFields *LineStringBuilder) {
	for _, field := range structure.Fields {
		accessModifier := "+"
		if unicode.IsLower(rune(field.Name[0])) {
			if !p.renderingOptions.PrivateMembers {
				continue
			}

			accessModifier = "-"
		}
		if accessModifier == "-" {
			privateFields.WriteLineWithDepth(2, fmt.Sprintf(`%s %s %s`, accessModifier, field.Name, field.Type))
		} else {
			publicFields.WriteLineWithDepth(2, fmt.Sprintf(`%s %s %s`, accessModifier, field.Name, field.Type))
		}
	}
}

// Returns an initialized struct of the given name or returns the existing one if it was already created
func (p *ClassParser) getOrCreateStruct(pkgName string, structName string) *Struct {
	result, ok := p.structure[pkgName][structName]
	if !ok {
		result = &Struct{
			PackageName:         pkgName,
			Functions:           make([]*Function, 0),
			Fields:              make([]*Field, 0),
			Type:                "",
			Composition:         make(map[string]struct{}, 0),
			Extends:             make(map[string]struct{}, 0),
			Aggregations:        make(map[string]struct{}, 0),
			PrivateAggregations: make(map[string]struct{}, 0),
		}
		p.structure[pkgName][structName] = result
	}
	return result
}

// Returns an existing struct only if it was created. nil otherwhise
func (p *ClassParser) getStruct(structName string) *Struct {
	split := strings.SplitN(structName, ".", 2)
	pack, ok := p.structure[split[0]]
	if !ok {
		return nil
	}
	return pack[split[1]]
}

// SetRenderingOptions Sets the rendering options for the Render() Function
func (p *ClassParser) SetRenderingOptions(ro map[RenderingOption]interface{}) error {
	for option, val := range ro {
		switch option {
		case RenderAggregations:
			p.renderingOptions.Aggregations = val.(bool)
		case RenderAliases:
			p.renderingOptions.Aliases = val.(bool)
		case RenderCompositions:
			p.renderingOptions.Compositions = val.(bool)
		case RenderFields:
			p.renderingOptions.Fields = val.(bool)
		case RenderImplementations:
			p.renderingOptions.Implementations = val.(bool)
		case RenderMethods:
			p.renderingOptions.Methods = val.(bool)
		case RenderConnectionLabels:
			p.renderingOptions.ConnectionLabels = val.(bool)
		case RenderTitle:
			p.renderingOptions.Title = val.(string)
		case RenderNotes:
			p.renderingOptions.Notes = val.(string)
		case AggregatePrivateMembers:
			p.renderingOptions.AggregatePrivateMembers = val.(bool)
		case RenderPrivateMembers:
			p.renderingOptions.PrivateMembers = val.(bool)
		case RenderTheme:
			p.renderingOptions.Theme = val.(string)
		default:
			return fmt.Errorf("Invalid Rendering option %v", option)
		}
	}
	return nil
}
func generateRenamedStructName(currentName string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(currentName, "")
}
