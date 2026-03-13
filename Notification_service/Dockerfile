FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o consumer .
 
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/consumer .

EXPOSE 8080

CMD ["./consumer"]

