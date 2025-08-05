# Hardlinker

A simple HTTP server that creates hard links between files.

## Features

- Runs an HTTP server on port 5070
- Accepts POST requests to `/link` with JSON body containing `source` and `destination`
- Creates a hard link from source to destination
- Creates the destination directory if needed
- Returns error if destination already exists
- Inserts the link into a YAML file (either cwd/links.yaml or LINK_YAML environment variable)

## API

### POST /link

Creates a hard link from source to destination.

**Request body:**
```json
{
  "source": "/path/to/source/file",
  "destination": "/path/to/destination/file"
}
```

**Response:**
- 200 OK on success
- 400 Bad Request if JSON is invalid or required fields missing
- 409 Conflict if destination already exists
- 500 Internal Server Error for other errors

## YAML File Format

The links are stored in a YAML file as an array of objects:

```yaml
- source: "/path/to/source/file"
  destination: "/path/to/destination/file"
```

## GitHub Actions

This project uses GitHub Actions for building and testing. The following workflows are available:

### Build Client
A workflow that builds the hardlinker binary for multiple platforms and uploads it as an artifact. This workflow runs automatically on semantic version tags (e.g., v1.0.0, v2.1.3).

The built binaries include:
- `hardlinker-linux-amd64` - Linux 64-bit binary
- `hardlinker-windows-amd64.exe` - Windows 64-bit binary
- `hardlinker-macos-amd64` - macOS 64-bit binary

These artifacts are automatically uploaded when a semver tag is pushed to the repository. Additionally, a GitHub release is created with these binaries attached.
The YAML file location can be specified via the `LINK_YAML` environment variable, or defaults to `cwd/links.yaml`.