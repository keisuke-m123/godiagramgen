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

type test struct {
	field  int
	field2 TestComplicatedAlias
	foo    parenthesizedtypedeclarations.Foo
}

type myInt int

var globalVariable int

//TestComplicatedAlias for testing purposes only
type TestComplicatedAlias func(strings.Builder) bool
