FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache git

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g cmd/app/main.go

RUN CGO_ENABLED=0 go build -o /app/server ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

COPY ./configs ./configs

COPY ./migrations ./migrations

EXPOSE 8080

CMD ["./server"]