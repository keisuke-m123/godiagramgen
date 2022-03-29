//lint:file-ignore U1000 Ignore all unused code, it's for testing support.
package testingsupport

import (
	f "fmt"
	"strings"

	"github.com/keisuke-m123/godiagramgen/testingsupport/parenthesizedtypedeclarations"
)

func (t *test) test() {
	f.Println("Hello Test")
}

type testInterface interface {
	test()
}

type test struct {
	field  int
	field2 TestComplicatedAlias
	foo    parenthesizedtypedeclarations.Foo
}

type definedTypeInt int

type aliasString = string

var globalVariable int

type definedTypeFunc func() *definedTypeInt

//TestComplicatedAlias for testing purposes only
type TestComplicatedAlias func(strings.Builder) bool
