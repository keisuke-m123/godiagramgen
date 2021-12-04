package parser

import (
	"fmt"
	"go/types"
	"strings"
)

const packageConstant = "{packageName}"

//Field can hold the name and type of any field
type Field struct {
	Name     string
	Type     string
	FullType string
}

//Returns a string representation of the given expression if it was recognized.
//Refer to the implementation to see the different string representations.
func getFieldType(typ types.Type, imports map[string]string) string {
	switch t := typ.(type) {
	case *types.Basic:
		return getBasic(t, imports)
	case *types.Slice:
		return getSliceType(t, imports)
	case *types.Array:
		return getArrayType(t, imports)
	case *types.Map:
		return getMapType(t, imports)
	case *types.Pointer:
		return getPointerType(t, imports)
	case *types.Chan:
		return getChanType(t, imports)
	case *types.Struct:
		return getStructType(t, imports)
	case *types.Interface:
		return getInterfaceType(t, imports)
	case *types.Signature:
		return getFuncType(t, imports)
	case *types.Named:
		return getNamedType(t, imports)
	case *types.Tuple:
		return ""
	default:
		return ""
	}
}

func getBasic(t *types.Basic, imports map[string]string) string {
	return t.Name()
}

func getSliceType(t *types.Slice, imports map[string]string) string {
	eType := getFieldType(t.Elem(), imports)
	return fmt.Sprintf("[]%s", eType)
}

func getArrayType(t *types.Array, imports map[string]string) string {
	eType := getFieldType(t.Elem(), imports)
	return fmt.Sprintf("[]%s", eType)
}

func getMapType(t *types.Map, imports map[string]string) string {
	kType := getFieldType(t.Key(), imports)
	eType := getFieldType(t.Elem(), imports)
	return fmt.Sprintf("map[%s]%s", kType, eType)
}

func getPointerType(t *types.Pointer, imports map[string]string) string {
	eType := getFieldType(t.Elem(), imports)
	return fmt.Sprintf("*%s", eType)
}

func getChanType(t *types.Chan, imports map[string]string) string {
	eType := getFieldType(t.Elem(), imports)
	return fmt.Sprintf("chan %s", eType)
}

func getStructType(t *types.Struct, imports map[string]string) string {
	fieldList := make([]string, 0)
	for i := 0; i < t.NumFields(); i++ {
		fType := getFieldType(t.Field(i).Type(), imports)
		fieldList = append(fieldList, fType)
	}
	return fmt.Sprintf("struct{%s}", strings.Join(fieldList, ", "))
}

func getInterfaceType(t *types.Interface, imports map[string]string) string {
	methods := make([]string, 0)
	for i := 0; i < t.NumMethods(); i++ {
		m := t.Method(i)
		methods = append(methods, fmt.Sprintf("%s %s", m.Name(), getFieldType(m.Type(), imports)))
	}
	return fmt.Sprintf("interface{%s}", strings.Join(methods, "; "))
}

func getFuncType(t *types.Signature, imports map[string]string) string {
	function := getFunction(t, "", imports, "")
	params := make([]string, 0)
	for _, pa := range function.Parameters {
		params = append(params, pa.Type)
	}

	returns := ""
	returnList := append([]string{}, function.ReturnValues...)
	if len(returnList) > 1 {
		returns = fmt.Sprintf("(%s)", strings.Join(returnList, ", "))
	} else {
		returns = strings.Join(returnList, "")
	}
	return fmt.Sprintf("func(%s) %s", strings.Join(params, ", "), returns)
}

func getNamedType(t *types.Named, imports map[string]string) string {
	if t.Obj().Pkg() == nil {
		return fmt.Sprintf("%s%s", packageConstant, t.Obj().Name())
	}

	pkgName := t.Obj().Pkg().Name()
	if _, ok := imports[pkgName]; !ok {
		return fmt.Sprintf("%s%s", packageConstant, t.Obj().Name())
	}

	return fmt.Sprintf("%s.%s", pkgName, t.Obj().Name())
}

func getFundamentalTypes(typ types.Type, imports map[string]string) []string {
	switch t := typ.(type) {
	case *types.Basic:
		return []string{}
	case *types.Slice:
		return getFundamentalTypes(t.Elem(), imports)
	case *types.Array:
		return getFundamentalTypes(t.Elem(), imports)
	case *types.Map:
		keyRes := getFundamentalTypes(t.Key(), imports)
		elmRes := getFundamentalTypes(t.Elem(), imports)
		return append(keyRes, elmRes...)
	case *types.Pointer:
		return getFundamentalTypes(t.Elem(), imports)
	case *types.Chan:
		return getFundamentalTypes(t.Elem(), imports)
	case *types.Struct:
		return []string{}
	case *types.Interface:
		return []string{}
	case *types.Signature:
		return []string{}
	case *types.Named:
		return []string{getNamedType(t, imports)}
	case *types.Tuple:
		return []string{}
	default:
		return []string{}
	}
}

var globalPrimitives = map[string]struct{}{
	"bool":        {},
	"string":      {},
	"int":         {},
	"int8":        {},
	"int16":       {},
	"int32":       {},
	"int64":       {},
	"uint":        {},
	"uint8":       {},
	"uint16":      {},
	"uint32":      {},
	"uint64":      {},
	"uintptr":     {},
	"byte":        {},
	"rune":        {},
	"float32":     {},
	"float64":     {},
	"complex64":   {},
	"complex128":  {},
	"error":       {},
	"*bool":       {},
	"*string":     {},
	"*int":        {},
	"*int8":       {},
	"*int16":      {},
	"*int32":      {},
	"*int64":      {},
	"*uint":       {},
	"*uint8":      {},
	"*uint16":     {},
	"*uint32":     {},
	"*uint64":     {},
	"*uintptr":    {},
	"*byte":       {},
	"*rune":       {},
	"*float32":    {},
	"*float64":    {},
	"*complex64":  {},
	"*complex128": {},
	"*error":      {},
}

func isPrimitiveString(t string) bool {
	_, ok := globalPrimitives[t]
	return ok
}

func replacePackageConstant(field, packageName string) string {
	if packageName != "" {
		packageName = fmt.Sprintf("%s.", packageName)
	}
	return strings.Replace(field, packageConstant, packageName, 1)
}
