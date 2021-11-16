package parser

import (
	"fmt"
	"go/token"
	"go/types"
	"testing"
)

type NoMatchField struct {
	types.Type
}

func TestGetFieldType(t *testing.T) {
	tt := []struct {
		Name           string
		ExpectedResult string
		InputType      types.Type
	}{
		{
			Name:           "TestTypesBasic",
			ExpectedResult: "int",
			InputType:      types.Typ[types.Int],
		},
		{
			Name:           "TestNotPrimitive",
			ExpectedResult: fmt.Sprintf("%s%s", packageConstant, "TestClass"),
			InputType: types.NewNamed(
				types.NewTypeName(token.NoPos, nil, "TestClass", nil),
				types.NewStruct(nil, nil),
				nil,
			),
		},
		{
			Name:           "TestTypesSlice",
			ExpectedResult: "[]int",
			InputType:      types.NewSlice(types.Typ[types.Int]),
		},
		{
			Name:           "TestTypesArray",
			ExpectedResult: "[]int",
			InputType:      types.NewArray(types.Typ[types.Int], 0),
		},
		{
			Name:           "TestTypesMap",
			ExpectedResult: "map[string]int",
			InputType:      types.NewMap(types.Typ[types.String], types.Typ[types.Int]),
		},
		{
			Name:           "TestTypesPointer",
			ExpectedResult: "*int",
			InputType:      types.NewPointer(types.Typ[types.Int]),
		},
		{
			Name:           "TestTypesChan",
			ExpectedResult: "chan int",
			InputType:      types.NewChan(types.SendOnly, types.Typ[types.Int]),
		},
		{
			Name:           "TestTypesStruct",
			ExpectedResult: "struct{int, string}",
			InputType: types.NewStruct(
				[]*types.Var{
					types.NewField(token.NoPos, nil, "int", types.Typ[types.Int], false),
					types.NewField(token.NoPos, nil, "string", types.Typ[types.String], false),
				},
				nil,
			),
		},
		{
			Name:           "TestTypesInterfaceType",
			ExpectedResult: "interface{Foo func(*FooComposed) *FooComposed}",
			InputType: types.NewInterfaceType(
				[]*types.Func{
					types.NewFunc(
						token.NoPos,
						nil,
						"Foo",
						types.NewSignature(
							nil,
							types.NewTuple(
								types.NewVar(
									token.NoPos,
									nil,
									"var1",
									types.NewPointer(
										types.NewNamed(
											types.NewTypeName(token.NoPos, nil, "FooComposed", nil),
											types.NewStruct(nil, nil),
											nil,
										),
									),
								),
							),
							types.NewTuple(
								types.NewVar(
									token.NoPos,
									nil,
									"",
									types.NewPointer(
										types.NewNamed(
											types.NewTypeName(token.NoPos, nil, "FooComposed", nil),
											types.NewStruct(nil, nil),
											nil,
										),
									),
								),
							),
							false,
						),
					),
				},
				nil,
			).Complete(),
		},
		{
			Name:           "TestTypesFuncWithOneResult",
			ExpectedResult: "func(*FooComposed) *FooComposed",
			InputType: types.NewSignature(
				nil,
				types.NewTuple(
					types.NewVar(
						token.NoPos,
						nil,
						"var1",
						types.NewPointer(
							types.NewNamed(
								types.NewTypeName(token.NoPos, nil, "FooComposed", nil),
								types.NewStruct(nil, nil),
								nil,
							),
						),
					),
				),
				types.NewTuple(
					types.NewVar(
						token.NoPos,
						nil,
						"",
						types.NewPointer(
							types.NewNamed(
								types.NewTypeName(token.NoPos, nil, "FooComposed", nil),
								types.NewStruct(nil, nil),
								nil,
							),
						),
					),
				),
				false,
			),
		},
		{
			Name:           "TestTypesFuncWithTwoResults",
			ExpectedResult: "func(*FooComposed) (*FooComposed, *string)",
			InputType: types.NewSignature(
				nil,
				types.NewTuple(
					types.NewVar(
						token.NoPos,
						nil,
						"var1",
						types.NewPointer(
							types.NewNamed(
								types.NewTypeName(token.NoPos, nil, "FooComposed", nil),
								types.NewStruct(nil, nil),
								nil,
							),
						),
					),
				),
				types.NewTuple(
					types.NewVar(
						token.NoPos,
						nil,
						"",
						types.NewPointer(
							types.NewNamed(
								types.NewTypeName(token.NoPos, nil, "FooComposed", nil),
								types.NewStruct(nil, nil),
								nil,
							),
						),
					),
					types.NewVar(token.NoPos, nil, "", types.NewPointer(types.Typ[types.String])),
				),
				false,
			),
		},
		{
			Name:           "TestNotMatchFieldType",
			ExpectedResult: "",
			InputType:      &NoMatchField{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			inputAliasMap := map[string]string{
				"puml": "goplantuml",
			}
			result := getFieldType(tc.InputType, inputAliasMap)
			if result != tc.ExpectedResult {
				t.Errorf("Expected result to be %s, got %s", tc.ExpectedResult, result)
			}
		})
	}
}

func TestIsPrimitiveStringPointer(t *testing.T) {
	if !isPrimitiveString("*int") {
		t.Errorf("TestIsPrimitiveStringPointer: expecting true, got false")
	}
}
