# GitHub Actions Workflows

## Docker Image Publishing

The `publish-docker.yml` workflow automatically builds and publishes Docker images to GitHub Container Registry (ghcr.io).

### Triggers

The workflow runs on:
- **Push to main branch**: Publishes with `latest` tag
- **Version tags** (e.g., `v1.0.0`, `v2.1.3`): Publishes with semantic version tags
- **Pull requests**: Builds but doesn't push (for validation)

### Image Tagging Strategy

Images are automatically tagged based on the trigger:

| Trigger | Tags Generated | Example |
|---------|---------------|---------|
| Push to main | `latest`, `main`, `main-<sha>` | `latest`, `main`, `main-abc1234` |
| Tag v1.2.3 | `1.2.3`, `1.2`, `1`, `v1.2.3` | All semantic versions |
| PR #42 | `pr-42` | For testing only (not pushed) |

### Using Published Images

After the workflow runs, your images will be available at:

```
ghcr.io/<username>/go-grpc
```

#### Pull the latest image:
```bash
docker pull ghcr.io/<username>/go-grpc:latest
```

#### Pull a specific version:
```bash
docker pull ghcr.io/<username>/go-grpc:1.2.3
```

#### Update docker-compose.yml to use published image:
```yaml
services:
  grpc-server:
    image: ghcr.io/<username>/go-grpc:${IMAGE_TAG:-latest}
    # Remove the build section when using published images
```

### Creating a Release

To publish a new version:

1. **Tag your release:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **The workflow automatically:**
   - Builds the Docker image
   - Tags it with semantic versions (1.0.0, 1.0, 1)
   - Pushes to GitHub Container Registry

### Permissions

The workflow uses `GITHUB_TOKEN` which is automatically provided by GitHub Actions. No additional secrets are required.

The image will be:
- **Public** - if your repository is public
- **Private** - if your repository is private

To make a private image public:
1. Go to your GitHub profile → Packages
2. Select your package
3. Package settings → Change visibility

### Build Cache

The workflow uses GitHub Actions cache to speed up builds:
- First build: ~3-5 minutes
- Subsequent builds: ~1-2 minutes (with cache)

### Monitoring

To view workflow runs:
1. Go to your repository on GitHub
2. Click the "Actions" tab
3. Select "Build and Publish Docker Image"

### Troubleshooting

**Image push fails:**
- Ensure GitHub Actions has package write permissions
- Check the Actions tab for detailed error logs

**Image not visible:**
- Check package visibility settings
- Verify you're logged in: `echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin`

**Build fails:**
- Check if `database.db` exists in the backend directory
- Review Dockerfile for any issues
- Check workflow logs in the Actions tab

