# Build stage
FROM golang:alpine AS build-stage

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o pdf-service main.go

# Production stage
FROM alpine:3.22.1

# Install Chrome and dependencies for headless PDF generation
RUN apk add --no-cache \
    chromium \
    fontconfig \
    ttf-dejavu \
    ttf-liberation \
    && rm -rf /var/cache/apk/*

# Set Chrome path
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Create temp directory and set permissions
RUN mkdir -p /tmp/pdf-conversion && \
    chown -R appuser:appgroup /tmp/pdf-conversion && \
    chmod 755 /tmp/pdf-conversion

# Copy the service binary from build stage
COPY --from=build-stage /app/pdf-service .

# Change ownership of the binary
RUN chown appuser:appgroup pdf-service

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8081

# Run the service
CMD ["./pdf-service"]
