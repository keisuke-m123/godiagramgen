package parser

import (
	"go/types"
	"reflect"
)

//Function holds the signature of a function with name, Parameters and Return values
type Function struct {
	Name                 string
	Parameters           []*Field
	ReturnValues         []string
	PackageName          string
	FullNameReturnValues []string
}

//SignaturesAreEqual Returns true if the two functions have the same signature (parameter names are not checked)
func (f *Function) SignaturesAreEqual(function *Function) bool {
	result := true
	result = result && (function.Name == f.Name)
	result = result && reflect.DeepEqual(f.FullNameReturnValues, function.FullNameReturnValues)
	result = result && (len(f.Parameters) == len(function.Parameters))
	if result {
		for i, p := range f.Parameters {
			if p.FullType != function.Parameters[i].FullType {
				return false
			}
		}
	}
	return result
}

// generate and return a function object from the given Functype. The names must be passed to this
// function since the FuncType does not have this information
func getFunction(s *types.Signature, name string, aliases map[string]string, packageName string) *Function {
	function := &Function{
		Name:                 name,
		Parameters:           make([]*Field, 0),
		ReturnValues:         make([]string, 0),
		FullNameReturnValues: make([]string, 0),
		PackageName:          packageName,
	}

	params := s.Params()
	if params != nil {
		for i := 0; i < params.Len(); i++ {
			pa := params.At(i)
			typ := getFieldType(pa.Type(), aliases)
			function.Parameters = append(function.Parameters, &Field{
				Name:     pa.Name(),
				Type:     replacePackageConstant(typ, ""),
				FullType: replacePackageConstant(typ, packageName),
			})
		}
	}

	results := s.Results()
	if results != nil {
		for i := 0; i < results.Len(); i++ {
			pa := results.At(i)
			typ := getFieldType(pa.Type(), aliases)
			function.ReturnValues = append(function.ReturnValues, replacePackageConstant(typ, ""))
			function.FullNameReturnValues = append(function.FullNameReturnValues, replacePackageConstant(typ, packageName))
		}
	}

	return function
}
