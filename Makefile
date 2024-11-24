HTMX_VERSION = 2.0.3

format::
	@go fmt -s -w .

build::
	@templ generate

run:: build
	DOMAIN=localhost go run main.go

update::
	@curl -s -o ./static/htmx.min.js https://unpkg.com/htmx.org@2.0.3/dist/htmx.min.js