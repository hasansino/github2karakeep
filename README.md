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
  ghcr.io/hasansino/github2karakeep:latest
```

docker-compose

```yaml
services:
  github2karakeep:
    image: ghcr.io/hasansino/github2karakeep:latest
    environment:
      TIMEOUT=10s
      GH_USERNAME={username}
      GH_TOKEN={token}
      GH_PER_PAGE=10
      KK_HOST={host}
      KK_TOKEN={token}
      KK_LIST=github2karakeep
      UPDATE_INTERVAL=24h
      EXPORT_LIMIT=10
```

## Configuration

cli arg | environment variable

### --timeout | TIMEOUT

Timeout for HTTP requests.
Duration format: "2h45m30s"
Default 10s.

### --gh-user | GH_USERNAME

GitHub username.
Required.

### --gh-token | GH_TOKEN

GitHub token. Token should have starring/read-only permission.
Required.

### --gh-per-page | GH_PER_PAGE

How many repos to fetch per page.
Default 10.

### --kk-host | KK_HOST

Karakeep host, including schema.
Required.

### --kk-token | KK_TOKEN

Karakeep API token.
Required.

### --kk-list | KK_LIST

Karakeep list name.
Default "github2karakeep".

### --update-interval | UPDATE_INTERVAL

Update interval.
Duration format: "2h45m30s"
Default 24h.

### --export-limit | EXPORT_LIMIT

Limit how many repos to export per run.
Default 10.