FROM golang:1.26-alpine AS builder
WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/svc ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache --upgrade ca-certificates tzdata curl
WORKDIR /app

COPY --from=builder /out/svc ./svc
COPY ./docs ./docs
COPY ./migrations ./migrations

CMD ["./svc"]
