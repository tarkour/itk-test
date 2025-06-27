ARG GO_VERSION=1.24.3
FROM golang:${GO_VERSION}-alpine
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ENTRYPOINT ["goose"]