GOPATH=$(shell go env GOPATH)

release:
	go build -o proxy main.go

debug:
	go run main.go -debug

prepare:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin v1.50.1
	git config core.hooksPath .githooks
	$(MAKE) lint

lint:
	${GOPATH}/bin/golangci-lint run
