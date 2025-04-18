# github2karakeep

Export starred repositories from GitHub to a KaraKeep list.

## Running

docker

```bash
docker run --rm \
  -e TIMEOUT=10s \
  -e GH_USERNAME={username} \
  -e GH_TOKEN={token} \
  -e GH_PER_PAGE=10 \
  -e KK_HOST={host} \
  -e KK_TOKEN={token} \
  -e KK_LIST=github2karakeep \
  -e UPDATE_INTERVAL=24h \
  -e EXPORT_LIMIT=10 \
  -e DEFAULT_TAG=github2karakeep \
  ghcr.io/hasansino/github2karakeep:latest
```

docker-compose

```yaml
services:
  github2karakeep:
    image: ghcr.io/hasansino/github2karakeep:latest
    environment:
      - TIMEOUT=10s
      - GH_USERNAME={username}
      - GH_TOKEN={token}
      - GH_PER_PAGE=10
      - KK_HOST={host}
      - KK_TOKEN={token}
      - KK_LIST=github2karakeep
      - UPDATE_INTERVAL=24h
      - EXPORT_LIMIT=10
      - DEFAULT_TAG=github2karakeep
```

## Configuration

| CLI Argument        | Environment Variable | Description                                                    | Default Value     |
|---------------------|----------------------|----------------------------------------------------------------|-------------------|
| `--timeout`         | `TIMEOUT`            | Timeout for HTTP requests. Duration format: `2h45m30s`.        | `10s`             |
| `--gh-user`         | `GH_USERNAME`        | GitHub username. **Required**.                                 |                   |
| `--gh-token`        | `GH_TOKEN`           | GitHub token with starring/read-only permission. **Required**. |                   |
| `--gh-per-page`     | `GH_PER_PAGE`        | Number of repositories to fetch per page.                      | `10`              |
| `--kk-host`         | `KK_HOST`            | KaraKeep host, including schema. **Required**.                 |                   |
| `--kk-token`        | `KK_TOKEN`           | KaraKeep API token. **Required**.                              |                   |
| `--kk-list`         | `KK_LIST`            | KaraKeep list name.                                            | `github2karakeep` |
| `--update-interval` | `UPDATE_INTERVAL`    | Update interval. Duration format: `2h45m30s`.                  | `24h`             |
| `--export-limit`    | `EXPORT_LIMIT`       | Limit the number of repositories to export per run.            | `10`              |
| `--default-tag`     | `DEFAULT_TAG`        | Default tag to add to every bookmark. Leave empty to omit.     | `github2karakeep` |

## Notes

+ Ensure your GitHub token has the necessary permissions (starring/read-only).
+ KaraKeep host should include the schema (e.g., https://example.com).
+ Use the UPDATE_INTERVAL to control how often the export process runs.
