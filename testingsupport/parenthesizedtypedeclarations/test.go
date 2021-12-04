//lint:file-ignore U1000 Ignore all unused code, it's generated
package parenthesizedtypedeclarations

type (
	//Foo is a test interface for testing purposes
	Foo interface {
		Foo()
	}
	//Bar is a test interface for testing purposes
	Bar interface {
		Bar()
	}

	defaultFoo struct{}
)

func (d *defaultFoo) Foo() {}
