FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o main .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
EXPOSE 8080
CMD ["./main"]
