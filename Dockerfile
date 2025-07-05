FROM golang:1.23-alpine AS builder

WORKDIR /app/twitter

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /app/prod_binary ./cmd/api/

FROM alpine:latest

ENV TZ=America/Argentina/Buenos_Aires

WORKDIR /app

COPY --from=builder /app/prod_binary /app/prod_binary

EXPOSE 8080

CMD ["/app/prod_binary"]
