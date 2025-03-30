# Build stage
FROM golang:1.24-alpine
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/uploads ./uploads

# Create directory for uploads
RUN mkdir -p /app/uploads && \
    chmod -R 777 /app/uploads

EXPOSE 8080

CMD ["./main"]