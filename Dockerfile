# Build the application from source
FROM golang:1.23.4 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /user-service -ldflags '-extldflags "-static"' ./cmd/user-service/main.go

# Run the tests in the container
FROM build-stage AS run-test-stage

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /user-service /user-service

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/user-service"]