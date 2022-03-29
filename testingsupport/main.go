//lint:file-ignore U1000 Ignore all unused code, it's for testing support.
package testingsupport

import (
	f "fmt"
	"strings"
	"time"

	"github.com/keisuke-m123/godiagramgen/testingsupport/parenthesizedtypedeclarations"
)

func (t *test) test() {
	f.Println("Hello Test")
}

type testInterface interface {
	test()
	returnTime() time.Time
}

type test struct {
	field  int
	field2 TestComplicatedAlias
	field3 time.Time
	foo    parenthesizedtypedeclarations.Foo
}

type definedTypeInt int

type definedTypeTime time.Time

type aliasString = string

var globalVariable int

type definedTypeFunc func() *definedTypeInt

//TestComplicatedAlias for testing purposes only
type TestComplicatedAlias func(strings.Builder) bool
