FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./src/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]
