FROM golang:1.24.6-alpine AS builder
WORKDIR /internal
RUN apk update  \
    && apk add ca-certificates --no-cache \
    && rm -rf /var/cache/apk/*
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN go build -mod=readonly -o build/irj-backend ./cmd/...

FROM scratch AS irj-backend
WORKDIR /opt/
COPY --from=builder /internal/build/irj-backend .

EXPOSE 5000
ENTRYPOINT ["/opt/irj-backend"]
