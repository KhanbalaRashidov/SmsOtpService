# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN swag init -g cmd/api/main.go ./docs

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget

RUN adduser -D -s /bin/sh appuser

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/docs ./docs

COPY --from=builder /app/.env* ./

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]