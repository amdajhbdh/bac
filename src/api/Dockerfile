FROM golang:alpine AS builder

WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary (adjust if your main is not at root, e.g. cmd/api/main.go)
RUN CGO_ENABLED=0 GOOS=linux go build -o /bac-api ./main.go   # or ./cmd/... if main is deeper

# Stage 2: Minimal runtime image
FROM alpine:3.21

# Install ca-certificates if your app makes HTTPS calls
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy only the binary from builder
COPY --from=builder /bac-api .

# Expose port (matches your compose ports mapping)
EXPOSE 8080

# Run the app
CMD ["/app/bac-api"]
