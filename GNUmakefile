default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

testacc:
	TF_ACC=1 go test -v -count=1 -cover -timeout 120m ./...

.PHONY: fmt lint testacc build install generate
