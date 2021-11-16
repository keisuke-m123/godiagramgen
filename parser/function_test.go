package parser

import (
	"go/token"
	"go/types"
	"reflect"
	"testing"
)

func TestGetFunction(t *testing.T) {
	tt := []struct {
		Name           string
		Func           *types.Signature
		ExpectedResult *Function
		FunctionName   string
	}{
		{
			Name: "Function with two typed parameters",
			Func: types.NewSignature(
				nil,
				types.NewTuple(
					types.NewVar(token.NoPos, nil, "param1", types.Typ[types.Int]),
					types.NewVar(token.NoPos, nil, "param2", types.Typ[types.Int]),
				),
				nil,
				false,
			),
			ExpectedResult: &Function{
				Name:        "TestFunction",
				PackageName: "main",
				Parameters: []*Field{
					{
						Name:     "param1",
						Type:     "int",
						FullType: "int",
					},
					{
						Name:     "param2",
						Type:     "int",
						FullType: "int",
					},
				},
				ReturnValues:         []string{},
				FullNameReturnValues: []string{},
			},
			FunctionName: "TestFunction",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {

			function := getFunction(tc.Func, tc.FunctionName, map[string]string{
				"main": "main",
			}, "main")

			if !reflect.DeepEqual(function, tc.ExpectedResult) {
				t.Errorf("Expected function to be %+v, got %+v", tc.ExpectedResult, function)
			}
		})
	}
}
