FROM rust:1.75-slim AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

COPY Cargo.toml Cargo.lock ./
COPY crates/ocr ./crates/ocr

RUN cargo build --release -p bac-ocr

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    libssl3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/target/release/bac-ocr /app/

EXPOSE 8082

CMD ["/app/bac-ocr"]
