FROM golang:1.24.0-alpine

WORKDIR /app

COPY . .

RUN go mod download

# サーバーコードをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o app server/main.go

EXPOSE 10000

CMD ["./app"] 