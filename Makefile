.PHONY: init
init:
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: install
install:
	go install ./cmd/godiagramgen

.PHONY: render
render:
	godiagramgen class --recursive --output=class-diagram.puml --theme=reddress-darkorange ./gocode
	godiagramgen package --output=./package-diagram.puml --theme=reddress-darkorange --ignore=./testingsupport .
	godiagramgen class --recursive --output=./testingsupport/testingsupport-all.puml --theme=reddress-darkorange ./testingsupport
	godiagramgen class --recursive --output=./testingsupport/testingsupport-all-ignore-directories.puml --ignore=./testingsupport/subfolder,./testingsupport/subfolder2,./testingsupport/connectionlabels ./testingsupport
	godiagramgen class --title='Test Title' --notes='Example 1,Example 1 continues,Example 2' --output=./testingsupport/testingsupport.puml ./testingsupport
	godiagramgen class --output=./testingsupport/testingsupport-parenthesizedtypedeclarations.puml ./testingsupport ./testingsupport/parenthesizedtypedeclarations
	godiagramgen class --output=./testingsupport/aliasmethods.puml ./testingsupport/aliasmethods
	godiagramgen class --output=./testingsupport/subfolder1-3.puml ./testingsupport/subfolder ./testingsupport/subfolder2 ./testingsupport/subfolder3
.PHONY: test
test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: check
check:
	goimports -w .
	go vet ./...
	staticcheck ./...

