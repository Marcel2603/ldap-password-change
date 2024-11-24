HTMX_VERSION = 2.0.3

format:
	@go fmt -s -w .

build: generate-static generate-dynamic
	@go build .

run: generate-dynamic
	HOST=localhost go run main.go

generate: generate-static generate-dynamic

generate-static:
	@mkdir -p ./static
	@curl -s -o ./static/htmx.min.js https://unpkg.com/htmx.org@${HTMX_VERSION}/dist/htmx.min.js

generate-dynamic:
	@go generate .
