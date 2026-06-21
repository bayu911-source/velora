# ==========================================
# STAGE 1: Build Stage
# ==========================================
FROM golang:1.22-alpine AS builder

# Install git dan certificates jika package Go membutuhkan unduhan dari repositori private/public
RUN apk add --no-cache git ca-certificates

# Tentukan working directory di dalam kontainer
WORKDIR /app

# Copy go.mod dan go.sum terlebih dahulu untuk memanfaatkan caching layer Docker
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code (termasuk main.go) ke dalam kontainer
COPY . .

# Compile aplikasi Go ke file binary tunggal bernama 'velora'
# Menargetkan langsung ke main.go untuk menghindari error multi-package
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/velora main.go

# ==========================================
# STAGE 2: Final Run Stage
# ==========================================
FROM alpine:3.19

# Install ca-certificates agar aplikasi bisa melakukan HTTPS request (misal ke Gemini API)
RUN apk add --no-cache ca-certificates tzdata

# Copy file binary yang sudah di-build dari stage 'builder'
COPY --from=builder /usr/local/bin/velora /usr/local/bin/velora

# Expose port sesuai yang digunakan aplikasi (8080)
EXPOSE 8080

# Jalankan binary secara default
ENTRYPOINT ["/usr/local/bin/velora"]
