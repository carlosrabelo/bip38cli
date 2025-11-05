# Docker Support

BIP38CLI bundles container artefacts under `docker/` so the repository root stays uncluttered.

## Quick Start

### Build the Docker image

```bash
docker build -f docker/Dockerfile -t bip38cli:latest .
```

### Run the container

```bash
# Show help
docker run --rm bip38cli:latest --help

# Encrypt a private key
docker run --rm -it bip38cli:latest encrypt 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ

# Decrypt a key
docker run --rm -it bip38cli:latest decrypt 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
```

### Convenience wrapper

For a CLI-like experience without installing Go, use the helper script:

```bash
./scripts/bip38cli-docker.sh --help
./scripts/bip38cli-docker.sh encrypt --verbose
```

The script builds an image tagged `bip38cli:local` (if missing or when `--build` is passed) and forwards every argument to the container.

## Docker Compose

### Production usage

```bash
# Start the service
docker compose -f docker/docker-compose.yml up bip38cli

# Run with custom command
docker compose -f docker/docker-compose.yml run --rm bip38cli encrypt --help
```

### Development

```bash
# Start development with hot reload
docker compose -f docker/docker-compose.yml --profile dev up bip38cli-dev

# Run tests
docker compose -f docker/docker-compose.yml --profile test up --build bip38cli-test
```

## Configuration

### Environment variables

- `BIP38CLI_CONFIG`: Path to configuration file (default: `/app/config/.bip38cli.yaml`)
- `VERSION`: Build version (for Docker build arg)
- `BUILD_TIME`: Build timestamp (for Docker build arg)

### Volumes

- `../config:/app/config`: Configuration files (read-only in production)
- `../data:/app/data`: Data directory for persistent storage

## Build Arguments

```bash
docker build \
  -f docker/Dockerfile \
  --build-arg VERSION=v1.0.0 \
  --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -t bip38cli:v1.0.0 .
```

## Multi-platform builds

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  -f docker/Dockerfile \
  -t bip38cli:latest \
  --push .
```

## Security Features

- **Non-root user**: Container runs as user `bip38cli` (UID 1001)
- **Minimal base image**: Uses Alpine Linux for reduced attack surface
- **Static binary**: No external dependencies required
- **Read-only filesystem**: Configuration files mounted as read-only

## Development with Hot Reload

For development, use the dev Dockerfile with Air hot reload:

```bash
docker compose -f docker/docker-compose.yml --profile dev up bip38cli-dev
```

This will:
- Watch for file changes
- Automatically rebuild and restart the application
- Provide live logging

## Production Deployment

### Docker Compose (Production)

```yaml
version: '3.8'
services:
  bip38cli:
    image: bip38cli:latest
    restart: unless-stopped
    user: "1001:1001"
    read_only: true
    tmpfs:
      - /tmp
    environment:
      - BIP38CLI_CONFIG=/app/config/.bip38cli.yaml
    volumes:
      - ./config:/app/config:ro
      - ./data:/app/data
```

### Kubernetes

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: bip38cli
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1001
    runAsGroup: 1001
    fsGroup: 1001
  containers:
  - name: bip38cli
    image: bip38cli:latest
    command: ["bip38cli", "--help"]
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop:
        - ALL
    volumeMounts:
    - name: config
      mountPath: /app/config
      readOnly: true
  volumes:
  - name: config
    configMap:
      name: bip38cli-config
```

## Troubleshooting

### Permission issues

If you encounter permission errors, ensure the data directory has correct permissions:

```bash
sudo chown -R 1001:1001 ./data
```

### Interactive mode

For interactive operations (like entering passphrases), use the `-it` flags:

```bash
docker run --rm -it bip38cli:latest encrypt
```

### Debug logging

Enable verbose logging:

```bash
docker run --rm bip38cli:latest --verbose encrypt --help
```
