.PHONY: init
init:
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: install
install:
	go install ./cmd/godiagramgen

.PHONY: render
render:
	godiagramgen class --recursive --show-connection-labels --show-aggregations --output=ParserDiagram.puml --theme=reddress-darkorange ./parser
	godiagramgen class --recursive --show-connection-labels --show-aggregations --output=TestingSupportDiagram.puml --theme=reddress-darkorange ./testingsupport
	godiagramgen class --title='Test Title' --notes='Example 1,Example 1 continues,Example 2' --output=./testingsupport/testingsupport.puml ./testingsupport
	godiagramgen class --show-aggregations --output=./testingsupport/testingsupport-parenthesizedtypedeclarations.puml ./testingsupport ./testingsupport/parenthesizedtypedeclarations
	godiagramgen class --output=./testingsupport/aliasmethods.puml ./testingsupport/aliasmethods
	godiagramgen class --hide-private-members --show-aggregations --output=./testingsupport/subfolder1-3.puml ./testingsupport/subfolder ./testingsupport/subfolder2 ./testingsupport/subfolder3


.PHONY: test
test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: check
check:
	goimports -w .
	go vet ./...
	staticcheck ./...

