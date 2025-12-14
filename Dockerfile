# ... (Stage Builder Tetap Sama) ...
FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod edit -go=1.23
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# ... (Stage Final) ...
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Jakarta

WORKDIR /root/

COPY --from=builder /app/server .

# PERBAIKAN: Copy folder seeds sesuai struktur aslinya
# Agar kode "pkg/database/seeds/data/*.json" bisa menemukan filenya
RUN mkdir -p pkg/database/seeds/data
COPY --from=builder /app/pkg/database/seeds/data ./pkg/database/seeds/data

ENV PORT=8086
EXPOSE 8086

CMD ["./server"]