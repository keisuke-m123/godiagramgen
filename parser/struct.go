package parser

import (
	"go/types"
)

//Struct represent a struct in golang, it can be of Type "class" or "interface" and can be associated
//with other structs via Composition and Extends
type Struct struct {
	PackageName         string
	Functions           []*Function
	Fields              []*Field
	Type                string
	Composition         map[string]struct{}
	Extends             map[string]struct{}
	Aggregations        map[string]struct{}
	PrivateAggregations map[string]struct{}
}

// ImplementsInterface returns true if the struct st conforms ot the given interface
func (st *Struct) ImplementsInterface(inter *Struct) bool {
	if len(inter.Functions) == 0 {
		return false
	}
	for _, f1 := range inter.Functions {
		foundMatch := false
		for _, f2 := range st.Functions {
			if f1.SignaturesAreEqual(f2) {
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			return false
		}
	}
	return true
}

//AddToComposition adds the composition relation to the structure. We want to make sure that *ExampleStruct
//gets added as ExampleStruct so that we can properly build the relation later to the
//class identifier
func (st *Struct) AddToComposition(fType string) {
	if len(fType) == 0 {
		return
	}
	if fType[0] == "*"[0] {
		fType = fType[1:]
	}
	st.Composition[fType] = struct{}{}
}

//AddToExtends Adds an extends relationship to this struct. We want to make sure that *ExampleStruct
//gets added as ExampleStruct so that we can properly build the relation later to the
//class identifier
func (st *Struct) AddToExtends(fType string) {
	if len(fType) == 0 {
		return
	}
	if fType[0] == "*"[0] {
		fType = fType[1:]
	}
	st.Extends[fType] = struct{}{}
}

//AddToAggregation adds an aggregation type to the list of aggregations
func (st *Struct) AddToAggregation(fType string) {
	st.Aggregations[fType] = struct{}{}
}

//addToPrivateAggregation adds an aggregation type to the list of aggregations for private members
func (st *Struct) addToPrivateAggregation(fType string) {
	st.PrivateAggregations[fType] = struct{}{}
}

//AddField adds a field into this structure. It parses the ast.Field and extract all
//needed information
func (st *Struct) AddField(field *types.Var, imports map[string]string) {
	typeName := getFieldType(field.Type(), imports)
	typeName = replacePackageConstant(typeName, "")
	if !field.Embedded() {
		typeName = replacePackageConstant(typeName, "")
		newField := &Field{
			Name: field.Name(),
			Type: typeName,
		}
		st.Fields = append(st.Fields, newField)
		for _, ft := range getFundamentalTypes(field.Type(), imports) {
			st.AddToAggregation(replacePackageConstant(ft, st.PackageName))
		}

	} else if field.Type() != nil {
		if typeName[0] == "*"[0] {
			typeName = typeName[1:]
		}
		st.AddToComposition(typeName)
	}
}

//AddMethod Parse the Field and if it is an ast.FuncType, then add the methods into the structure
func (st *Struct) AddMethod(method *types.Func, aliases map[string]string) {
	s, ok := method.Type().(*types.Signature)
	if !ok {
		return
	}
	function := getFunction(s, method.Name(), aliases, st.PackageName)
	st.Functions = append(st.Functions, function)
}
