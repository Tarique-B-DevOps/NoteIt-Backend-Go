# Build Stage
FROM golang:1.22.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final Stage
FROM alpine:latest  

WORKDIR /root/
RUN apk --no-cache add curl
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
