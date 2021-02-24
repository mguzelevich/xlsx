FROM golang:1.20-alpine AS builder


RUN apk add --no-cache git ca-certificates

COPY . /src/xlsx

WORKDIR /src/xlsx

ARG GO_IMPORT_TOKEN

RUN go env -w GOPRIVATE=github.com/mguzelevich/*

RUN git config --global url."https://${GO_IMPORT_TOKEN}@github.com/mguzelevich/".insteadOf "https://github.com/mguzelevich/"

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ../.build/xlsx ./cmd/xlsxcli

FROM alpine:3

COPY --from=builder /src/.build/xlsx /

# ENTRYPOINT [ "/xlsx" ]
