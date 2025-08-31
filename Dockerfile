# ---- Builder stage ----
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# ---- Run stage ----
FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8000

COPY .env .

CMD ["./main"]
