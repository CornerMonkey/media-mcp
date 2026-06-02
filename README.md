# media-mcp

Go MCP server for read-only media automation integrations.

## Current Scope

`media-mcp` currently provides read-only Sonarr and Radarr access patterns for:

- system status
- health checks
- root folders
- quality profiles
- series and movie lookup
- series and movie listing

Mutating operations are not currently implemented. `MEDIA_MCP_ALLOW_MUTATIONS` is parsed as a safety gate for future write actions and should remain unset or `false` for read-only use.

## Configuration

Set at least one complete service configuration:

| Variable | Description |
| --- | --- |
| `SONARR_URL` | Base Sonarr URL, for example `http://sonarr:8989` |
| `SONARR_API_KEY` | Sonarr API key |
| `RADARR_URL` | Base Radarr URL, for example `http://radarr:7878` |
| `RADARR_API_KEY` | Radarr API key |
| `MEDIA_MCP_ALLOW_MUTATIONS` | Set to `true` only when mutation tools exist and are intentionally enabled |

## Local Tests

```sh
GOCACHE=/private/tmp/media-mcp-go-build-cache go test ./...
```

## Docker

Build:

```sh
docker build -t media-mcp .
```

Run:

```sh
docker run --rm \
  -e SONARR_URL=http://sonarr:8989 \
  -e SONARR_API_KEY=replace-me \
  -e RADARR_URL=http://radarr:7878 \
  -e RADARR_API_KEY=replace-me \
  -e MEDIA_MCP_ALLOW_MUTATIONS=false \
  media-mcp
```
