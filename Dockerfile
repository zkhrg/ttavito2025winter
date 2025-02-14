# Build the application from source
FROM golang:1.23.4 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o /ttavito -ldflags '-extldflags "-static"' ./cmd/main.go

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage

WORKDIR /

COPY --from=build-stage /ttavito /ttavito

EXPOSE 8080

ENTRYPOINT ["/ttavito"]