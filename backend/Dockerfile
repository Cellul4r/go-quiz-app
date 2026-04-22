# Stage 1 — build
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./app/main.go

# Stage 2 — final image
FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]