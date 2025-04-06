FROM golang:1.24.0-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main server/main.go

EXPOSE 10000

CMD ["./main"] 