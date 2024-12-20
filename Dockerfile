# Fetch
FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN apk update && apk add make curl
RUN make generate
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app

# Deploy
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /app/static /static
COPY --from=build /app/bin/app /app
EXPOSE 3333
USER nonroot:nonroot
ENTRYPOINT ["/app"]