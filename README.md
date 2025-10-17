# doc-gen

A lightweight PDF generation service that converts HTML to PDF using headless Chromium. This service is containerized and ready for deployment.

## Features

- HTML to PDF conversion using headless Chromium
- RESTful API with health check endpoint
- Dockerized for easy deployment
- Lightweight Alpine Linux base image
- Non-root user for security

## API Endpoints

### Health Check
```
GET /health
```
Returns the service health status.

### Convert HTML to PDF
```
POST /convert
Content-Type: application/json

{
  "html": "<html><body><h1>Hello World</h1></body></html>",
  "title": "My Document"
}
```

## Docker Usage

### Using the Published Image

The service is published to GitHub Container Registry. You can use it in your projects:

```bash
# Pull the latest image
docker pull ghcr.io/YOUR_USERNAME/doc-gen:latest

# Run the service
docker run -p 8081:8081 ghcr.io/YOUR_USERNAME/doc-gen:latest
```

### Using in Docker Compose

```yaml
version: '3.8'
services:
  doc-gen:
    image: ghcr.io/YOUR_USERNAME/doc-gen:latest
    ports:
      - "8081:8081"
    environment:
      - CHROME_BIN=/usr/bin/chromium-browser
```

### Using in Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: doc-gen
spec:
  replicas: 1
  selector:
    matchLabels:
      app: doc-gen
  template:
    metadata:
      labels:
        app: doc-gen
    spec:
      containers:
      - name: doc-gen
        image: ghcr.io/YOUR_USERNAME/doc-gen:latest
        ports:
        - containerPort: 8081
---
apiVersion: v1
kind: Service
metadata:
  name: doc-gen-service
spec:
  selector:
    app: doc-gen
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8081
  type: LoadBalancer
```

## Local Development

### Building the Docker Image

```bash
# Build the image
docker build -t doc-gen .

# Run locally
docker run -p 8081:8081 doc-gen
```

### Testing the Service

```bash
# Health check
curl http://localhost:8081/health

# Convert HTML to PDF
curl -X POST http://localhost:8081/convert \
  -H "Content-Type: application/json" \
  -d '{"html": "<html><body><h1>Test PDF</h1><p>This is a test document.</p></body></html>", "title": "Test Document"}' \
  --output test.pdf
```

## Publishing to GitHub Container Registry

This repository includes a GitHub Actions workflow that automatically builds and pushes the Docker image to GHCR when you:

1. Push to the main branch
2. Create a tag (e.g., `v1.0.0`)
3. Open a pull request

### Setting Up Automated Publishing

1. Push this repository to GitHub
2. The GitHub Actions workflow will automatically run
3. Your image will be available at `ghcr.io/YOUR_USERNAME/doc-gen:latest`

### Version Tags

- `latest` - Always points to the latest version
- `v1.0.0` - Specific version tags
- `main` - Branch-based tags for development

## Configuration

The service runs on port 8081 by default. You can customize the Chromium settings by modifying the environment variables in the Dockerfile.

## Security

The Docker image runs as a non-root user for enhanced security. The service uses Alpine Linux for a minimal attack surface.
