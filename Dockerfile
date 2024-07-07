# Gunakan image dasar untuk Go
FROM golang:1.23.3-alpine

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum ke container
COPY go.mod go.sum ./

# Install dependensi
RUN go mod download

# Copy semua file proyek ke container
COPY . .

# Build aplikasi
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
