# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /agent ./cmd/main.go

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache \
    tesseract-ocr \
    tesseract-ocr-data-fra \
    tesseract-ocr-data-eng \
    poppler-utils \
    curl \
    bash

COPY --from=builder /agent /usr/local/bin/

RUN adduser -D -u 1000 bacuser

WORKDIR /home/bacuser

USER bacuser

ENTRYPOINT ["/usr/local/bin/agent"]
CMD ["--help"]
