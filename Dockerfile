FROM golang:1.24.2-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main server/main.go

EXPOSE 8080

CMD ["./main"] 