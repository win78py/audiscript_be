FROM golang:1.24 as builder

WORKDIR /app
COPY . .
RUN go build -o app ./cmd/api/main.go

FROM debian:bookworm-slim

# Cài Python3 và pip
RUN apt-get update && apt-get install -y python3 python3-pip

WORKDIR /app
COPY --from=builder /app/app .
COPY internal/transcribe/transcribe.py internal/transcribe/transcribe.py

# Copy các file cần thiết khác (nếu có)

EXPOSE 8080

CMD ["./app"]