.PHONY: install
install:
	go install ./cmd/godiagramgen

.PHONY: renderParser
renderParser:
	godiagramgen class --recursive --show-connection-labels --show-aggregations --output=ParserDiagram.puml --theme=reddress-darkorange ./parser

.PHONY: renderTestingSupport
renderTestingSupport:
	godiagramgen class --recursive --show-connection-labels --show-aggregations --output=TestingSupportDiagram.puml --theme=reddress-darkorange ./testingsupport