FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o main main.go
FROM alpine:latest
RUN apk add --no-cache sqlite
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/models/negri.db ./negri.db
EXPOSE 1488
CMD ["./main"]
