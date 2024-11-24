# Fetch
FROM golang:latest AS fetch-stage
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download

# Generate
FROM ghcr.io/a-h/templ:latest AS generate-stage
COPY . /app
WORKDIR /app
RUN ["templ", "generate"]

# Build
FROM golang:latest AS build-stage
COPY --from=generate-stage /app /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app

# Deploy
#FROM gcr.io/distroless/base-debian12 AS deploy-stage
#WORKDIR /
#COPY --from=build-stage /app/app /app
#EXPOSE 3333
#USER nonroot:nonroot
#ENTRYPOINT ["/app"]