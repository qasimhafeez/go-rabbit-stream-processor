FROM golang:1.22 AS builder

WORKDIR app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/gomod-cache \
    go mod download

COPY . ./
RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
    go build -o /main

CMD ["/main"]
