PKGS        := $(shell go list ./...)
HTMX_VERSION = 2.0.4
HTMX_RESPONSE_TARGET_VERSION = 2.0.3
BOOTSTRAP_VERSION = 5.3.3
ALPINE_JS_VERSION = 3.14.9

COVERMODE   := atomic
COVERFILE   := coverage.out
HTMLFILE    := coverage.html
GOFLAGS     := -race -shuffle=on -tags=test -covermode=$(COVERMODE) -coverprofile=$(COVERFILE)

format:
	@gofmt -l -s -w .

lint: format
	@go run github.com/mgechev/revive@latest -config config.toml -formatter friendly ./...

init-precommit:
	@pre-commit install


test:
	go test $(GOFLAGS) $(PKGS) -cover

cover:
	go test $(GOFLAGS) $(PKGS)
	@echo
	@go tool cover -func=$(COVERFILE) | tail -n1

cover-html: cover
	go tool cover -html=$(COVERFILE) -o $(HTMLFILE)
	@echo "Wrote $(HTMLFILE)"

build: generate-static generate-dynamic
	@go build -v -o bin .

run: generate-dynamic
	SERVER_HOST=localhost SERVER_PORT=4000 go run main.go

generate: generate-static generate-dynamic

generate-static:
	@mkdir -p ./static
	@curl -s -o ./static/htmx.min.js https://unpkg.com/htmx.org@${HTMX_VERSION}/dist/htmx.min.js
	@curl -s -o ./static/htmx-response-target.js https://unpkg.com/htmx-ext-response-targets@${HTMX_RESPONSE_TARGET_VERSION}/response-targets.js
	@curl -s -o ./static/bootstrap.min.css https://unpkg.com/bootstrap@${BOOTSTRAP_VERSION}/dist/css/bootstrap.min.css
	@curl -s -o ./static/bootstrap.min.css.map https://unpkg.com/bootstrap@${BOOTSTRAP_VERSION}/dist/css/bootstrap.min.css.map
	@curl -s -o ./static/bootstrap.bundle.min.js https://cdn.jsdelivr.net/npm/bootstrap@${BOOTSTRAP_VERSION}/dist/js/bootstrap.bundle.min.js
	@curl -s -o ./static/alpine.min.js https://unpkg.com/alpinejs@${ALPINE_JS_VERSION}/dist/cdn.min.js

generate-dynamic:
	@go generate .
