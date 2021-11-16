package parser

import (
	"go/token"
	"go/types"
	"reflect"
	"testing"
)

func TestStructImplementsInterface(t *testing.T) {
	tt := []struct {
		name           string
		structure      *Struct
		inter          *Struct
		expectedResult bool
	}{
		{
			name: "Correct implementation",
			structure: &Struct{
				Functions: []*Function{
					{
						Name: "foo",
						Parameters: []*Field{
							{
								Name:     "a",
								Type:     "int",
								FullType: "int",
							},
							{
								Name:     "b",
								Type:     "string",
								FullType: "string",
							},
						},
						FullNameReturnValues: []string{"int", "error"},
					},
				},
				Type: "class",
			},
			inter: &Struct{
				Functions: []*Function{
					{
						Name: "foo",
						Parameters: []*Field{
							{
								Type:     "int",
								FullType: "int",
							},
							{
								Type:     "string",
								FullType: "string",
							},
						},
						FullNameReturnValues: []string{"int", "error"},
					},
				},
				Type: "class",
			},
			expectedResult: true,
		}, {
			name: "Parameters not in order",
			structure: &Struct{
				Functions: []*Function{
					{
						Name: "foo",
						Parameters: []*Field{
							{
								Name:     "a",
								Type:     "int",
								FullType: "int",
							},
							{
								Name:     "b",
								Type:     "string",
								FullType: "string",
							},
						},
						FullNameReturnValues: []string{"int", "error"},
					},
				},
				Type: "interface",
			},
			inter: &Struct{
				Functions: []*Function{
					{
						Name: "foo",
						Parameters: []*Field{
							{
								Name:     "b",
								Type:     "string",
								FullType: "string",
							},
							{
								Name:     "a",
								Type:     "int",
								FullType: "int",
							},
						},
						FullNameReturnValues: []string{"int", "error"},
					},
				},
				Type: "interface",
			},
			expectedResult: false,
		}, {
			name: "Empty Interface",
			structure: &Struct{
				Functions: []*Function{
					{
						Name: "foo",
						Parameters: []*Field{
							{
								Name:     "a",
								Type:     "int",
								FullType: "int",
							},
							{
								Name:     "b",
								Type:     "string",
								FullType: "string",
							},
						},
						FullNameReturnValues: []string{"int", "error"},
					},
				},
				Type: "class",
			},
			inter: &Struct{
				Functions: []*Function{},
				Type:      "interface",
			},
			expectedResult: false,
		}, {
			name: "Different Function Names",
			structure: &Struct{
				Functions: []*Function{
					{
						Name: "foo",
						Parameters: []*Field{
							{
								Name:     "a",
								Type:     "int",
								FullType: "int",
							},
							{
								Name:     "b",
								Type:     "string",
								FullType: "String",
							},
						},
						FullNameReturnValues: []string{"int", "error"},
					},
				},
				Type: "class",
			},
			inter: &Struct{
				Functions: []*Function{
					{
						Name: "bar",
						Parameters: []*Field{
							{
								Name:     "a",
								Type:     "int",
								FullType: "int",
							},
							{
								Name:     "b",
								Type:     "string",
								FullType: "string",
							},
						},
						FullNameReturnValues: []string{"int", "error"},
					},
				},
				Type: "class",
			},
			expectedResult: false,
		}, {
			name: "Return value different",
			structure: &Struct{
				Functions: []*Function{
					{
						Name: "foo",
						Parameters: []*Field{
							{
								Name:     "a",
								Type:     "int",
								FullType: "int",
							},
							{
								Name:     "b",
								Type:     "string",
								FullType: "string",
							},
						},
						FullNameReturnValues: []string{"int", "error"},
					},
				},
				Type: "class",
			},
			inter: &Struct{
				Functions: []*Function{
					{
						Name: "foo",
						Parameters: []*Field{
							{
								Type: "int",
							},
							{
								Type: "string",
							},
						},
						FullNameReturnValues: []string{"error", "int"},
					},
				},
				Type: "class",
			},
			expectedResult: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.structure.ImplementsInterface(tc.inter)
			if result != tc.expectedResult {
				t.Errorf("Expected result to be %t, got %t", tc.expectedResult, result)
			}
		})

	}
}

func TestAddToComposition(t *testing.T) {
	st := &Struct{
		Functions: []*Function{
			{
				Name: "foo",
				Parameters: []*Field{
					{
						Type: "int",
					},
					{
						Type: "string",
					},
				},
				ReturnValues:         []string{"error", "int"},
				FullNameReturnValues: []string{"error", "int"},
			},
		},
		Type:        "class",
		PackageName: "test",
		Fields:      make([]*Field, 0),
		Composition: make(map[string]struct{}),
		Extends:     make(map[string]struct{}),
	}
	st.AddToComposition("Foo")

	if !arrayContains(st.Composition, "Foo") {
		t.Errorf("TestAddToComposition: Expected CompositionArray to have %s, but it contains %v", "Foo", st.Composition)
	}

	st.AddToComposition("")

	if arrayContains(st.Composition, "") {
		t.Errorf(`TestAddToComposition: Expected CompositionArray to not have "", but it contains %v`, st.Composition)
	}
	testArray := map[string]struct{}{
		"Foo": {},
	}
	if !reflect.DeepEqual(st.Composition, testArray) {

		t.Errorf("TestAddToComposition: Expected CompositionArray to be %v, but it contains %v", testArray, st.Composition)
	}

	st.AddToComposition("*Foo2")

	if !arrayContains(st.Composition, "Foo2") {
		t.Errorf("TestAddToComposition: Expected CompositionArray to have %s, but it contains %v", "Foo2", st.Composition)
	}
}
func TestAddToExtension(t *testing.T) {
	st := &Struct{
		Functions: []*Function{
			{
				Name: "foo",
				Parameters: []*Field{
					{
						Type: "int",
					},
					{
						Type: "string",
					},
				},
				ReturnValues:         []string{"error", "int"},
				FullNameReturnValues: []string{"error", "int"},
			},
		},
		Type:        "class",
		PackageName: "test",
		Fields:      make([]*Field, 0),
		Composition: make(map[string]struct{}),
		Extends:     make(map[string]struct{}),
	}
	st.AddToExtends("Foo")

	if !arrayContains(st.Extends, "Foo") {
		t.Errorf("TestAddToComposition: Expected Extends Array to have %s, but it contains %v", "Foo", st.Composition)
	}

	st.AddToExtends("")

	if arrayContains(st.Extends, "") {
		t.Errorf(`TestAddToComposition: Expected Extends Array to not have "", but it contains %v`, st.Composition)
	}
	testArray := map[string]struct{}{
		"Foo": {},
	}
	if !reflect.DeepEqual(st.Extends, testArray) {
		t.Errorf("TestAddToComposition: Expected Extends Array to be %v, but it contains %v", testArray, st.Composition)
	}

	st.AddToExtends("*Foo2")

	if !arrayContains(st.Extends, "Foo2") {
		t.Errorf("TestAddToComposition: Expected Extends Array to have %s, but it contains %v", "Foo2", st.Composition)
	}
}

func arrayContains(a map[string]struct{}, text string) bool {

	found := false
	for v := range a {
		if v == text {
			found = true
			break
		}
	}
	return found
}

func TestAddField(t *testing.T) {
	st := &Struct{
		PackageName: "main",
		Functions: []*Function{
			{
				Name:                 "foo",
				Parameters:           []*Field{},
				ReturnValues:         []string{"error", "int"},
				FullNameReturnValues: []string{"error", "int"},
			},
		},
		Type:         "class",
		Fields:       make([]*Field, 0),
		Composition:  make(map[string]struct{}),
		Extends:      make(map[string]struct{}),
		Aggregations: make(map[string]struct{}),
	}
	st.AddField(types.NewField(token.NoPos, nil, "foo", types.Typ[types.Int], false), make(map[string]string))
	if len(st.Fields) != 1 {
		t.Errorf("TestAddField: Expected st.Fields to have exactly one element but it has %d elements", len(st.Fields))
	}
	testField := &Field{
		Name: "foo",
		Type: "int",
	}
	if !reflect.DeepEqual(st.Fields[0], testField) {
		t.Errorf("TestAddField: Expected st.Fields[0] to have %v, got %v", testField, st.Fields[0])
	}
	st.AddField(types.NewField(
		token.NoPos,
		nil,
		"",
		types.NewPointer(types.NewNamed(
			types.NewTypeName(token.NoPos, nil, "FooComposed", nil),
			types.NewStruct(nil, nil),
			nil,
		)),
		true,
	), make(map[string]string))
	if !arrayContains(st.Composition, "FooComposed") {
		t.Errorf("TestAddField: Expecting FooComposed to be part of the compositions ,but the array had %v", st.Composition)
	}
	st.AddField(types.NewField(
		token.NoPos,
		nil,
		"Foo",
		types.NewPointer(types.NewNamed(
			types.NewTypeName(token.NoPos, nil, "FooComposed", nil),
			types.NewStruct(nil, nil),
			nil,
		)),
		false,
	), make(map[string]string))
	if !arrayContains(st.Aggregations, "main.FooComposed") {
		t.Errorf("TestAddField: Expecting main.FooComposed to be part of the aggregations ,but the array had %v", st.Aggregations)
	}
}

func TestAddMethod(t *testing.T) {
	st := &Struct{
		PackageName: "main",
		Functions:   []*Function{},
		Type:        "class",
	}
	st.AddMethod(types.NewFunc(token.NoPos, nil, "foo", nil), make(map[string]string))
	if len(st.Functions) != 0 {
		t.Errorf("TestAddMethod: Expected Functions array to be empty but it contains %v", st.Functions)
	}
	st.AddMethod(types.NewFunc(
		token.NoPos,
		nil,
		"foo",
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
	), make(map[string]string))
	if len(st.Functions) != 1 {
		t.Errorf("TestAddMethod: Expected st.Functions to have exactly one element but it has %d elements", len(st.Functions))
	}
	testFunction := &Function{
		PackageName: "main",
		Name:        "foo",
		Parameters: []*Field{
			{
				Name:     "var1",
				Type:     "*FooComposed",
				FullType: "*main.FooComposed",
			},
		},
		ReturnValues:         []string{"*FooComposed"},
		FullNameReturnValues: []string{"*main.FooComposed"},
	}
	if !st.Functions[0].SignaturesAreEqual(testFunction) {
		t.Errorf("TestAddMethod: Expected st.Function[0] to have %v, got %v", testFunction, st.Functions[0])
	}
}
