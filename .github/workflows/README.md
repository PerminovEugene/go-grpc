# GitHub Actions Workflows

## Docker Publish Workflow

The `docker-publish.yml` workflow automatically builds and publishes Docker images to GitHub Container Registry (GHCR) when version tags are pushed.

### Trigger

The workflow **only** runs when version tags matching the pattern `v*.*.*` are pushed to the repository (e.g., `v1.0.0`, `v2.1.3`).

### What it does

1. Checks out the repository
2. Logs in to GitHub Container Registry using the built-in `GITHUB_TOKEN`
3. Extracts version information from the tag
4. Builds the Docker image from `backend/Dockerfile`
5. Pushes the image with multiple tags:
   - Full semantic version (e.g., `v1.2.3`)
   - Major.minor version (e.g., `v1.2`)
   - Major version (e.g., `v1`)

### Creating a Release

To trigger a build and publish a new Docker image:

```bash
# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0
```

### Published Images

Images are available at:
```
ghcr.io/perminovEugene/go-grpc:v1.0.0
ghcr.io/perminovEugene/go-grpc:v1.0
ghcr.io/perminovEugene/go-grpc:v1
```

### Permissions

The workflow uses the built-in `GITHUB_TOKEN` which has the necessary permissions to:
- Read repository contents
- Write to GitHub Packages

No additional secrets configuration is required.

### Pulling Images

To use the published images:

```bash
# Pull a specific version
docker pull ghcr.io/perminovEugene/go-grpc:v1.0.0

# Pull the latest minor version
docker pull ghcr.io/perminovEugene/go-grpc:v1.0

# Pull the latest major version
docker pull ghcr.io/perminovEugene/go-grpc:v1
```

**Note:** Images in GHCR may be private by default. To make them public, go to the package settings on GitHub.
