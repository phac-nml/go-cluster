# syntax=docker/dockerfile:1

FROM golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY TestInputs/ ./TestInputs/

RUN CGO_ENABLED=0 GOOS=linux go install

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...