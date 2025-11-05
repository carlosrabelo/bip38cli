# Multi-stage build for Go 1.24+ with optimizations
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY core/go.mod core/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY core/ .

# Build the application with optimizations
ARG VERSION=dev
ARG BUILD_TIME=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -tags netgo,osusergo,static_build \
    -a -installsuffix cgo \
    -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}" \
    -o bip38cli ./cmd/bip38cli

# Final minimal image
FROM alpine:latest

ARG VERSION=dev

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S bip38cli && \
    adduser -u 1001 -S bip38cli -G bip38cli

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bip38cli /usr/local/bin/bip38cli

# Change ownership to non-root user
RUN chown bip38cli:bip38cli /usr/local/bin/bip38cli

# Switch to non-root user
USER bip38cli

# Set entrypoint
ENTRYPOINT ["bip38cli"]

# Default command shows help
CMD ["--help"]

# Labels for metadata
LABEL maintainer="carlosrabelo" \
      version="${VERSION}" \
      description="BIP38 Bitcoin private key encryption CLI tool" \
      org.opencontainers.image.title="bip38cli" \
      org.opencontainers.image.description="A CLI tool for BIP38 Bitcoin private key encryption" \
      org.opencontainers.image.vendor="bip38cli" \
      org.opencontainers.image.licenses="MIT"
