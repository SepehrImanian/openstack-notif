FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vm-monitor ./cmd/vm-monitor

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/vm-monitor /app/vm-monitor
USER nonroot:nonroot
ENTRYPOINT ["/app/vm-monitor"]