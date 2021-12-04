.PHONY: install
install:
	go install ./cmd/godiagramgen

.PHONY: render
render:
	godiagramgen class --recursive --show-connection-labels --show-aggregations --output=ParserDiagram.puml --theme=reddress-darkorange ./parser
	godiagramgen class --recursive --show-connection-labels --show-aggregations --output=TestingSupportDiagram.puml --theme=reddress-darkorange ./testingsupport

.PHONY: test
test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./parser