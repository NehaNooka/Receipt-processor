FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /Receipt-processor

FROM alpine:latest

WORKDIR /

COPY --from=builder /Receipt-processor /Receipt-processor

EXPOSE 8080

ENTRYPOINT ["/Receipt-processor"]