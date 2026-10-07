# Stage 1: Build
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy modul dan download dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh kode sumber
COPY . .

# Build aplikasi menjadi file binary bernama "main"
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Copy binary dari stage 1
COPY --from=builder /app/main .

# Buat folder uploads untuk berjaga-jaga jika volume gagal
RUN mkdir -p uploads

# Expose port
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]