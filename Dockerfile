FROM golang:1.24.0-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o task-recommender ./cmd/todo

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/task-recommender .

EXPOSE 8080
CMD ["./task-recommender", "serve"]
