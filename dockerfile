FROM golang:1.24 as builder

WORKDIR /app
COPY . .
RUN go build -o app ./cmd/api/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y python3 python3-pip ffmpeg

# Kiểm tra version
RUN python3 --version && pip3 --version

# Cài các thư viện Python cần thiết cho Whisper
RUN pip3 install --no-cache-dir openai-whisper torch requests

# Kiểm tra đã cài whisper chưa
RUN python3 -c "import whisper; print('Whisper installed OK')"

WORKDIR /app
COPY --from=builder /app/app .
COPY internal/transcribe/transcribe.py internal/transcribe/transcribe.py

EXPOSE 8080

CMD ["./app"]