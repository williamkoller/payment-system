FROM golang:1.24.5-alpine AS builder

WORKDIR /app

RUN mkdir -p /app

RUN apk add --no-cache git ca-certificates


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
