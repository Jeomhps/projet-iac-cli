# projet-iac-cli

Go CLI for the Projet IAC API at `https://localhost/api` (by default).

## Build

```bash
go mod tidy
go build -o projet-iac-cli
# optional: embed version/commit
# go build -ldflags "-X 'github.com/Jeomhps/projet-iac-cli/cmd.commit=$(git rev-parse --short HEAD)'" -o projet-iac-cli
```

## Quick start (dev)

```bash
export API_BASE=https://localhost
export API_PREFIX=/api
export VERIFY_TLS=false
export ADMIN_DEFAULT_USERNAME=admin
export ADMIN_DEFAULT_PASSWORD=change-me

./projet-iac-cli login
./projet-iac-cli whoami
./projet-iac-cli machines list
./projet-iac-cli machines add --name alpine-1 --host localhost --port 22221 --user root --password test
./projet-iac-cli reservations
./projet-iac-cli reserve --count 2 --duration 60 --password test
./projet-iac-cli release-all
./projet-iac-cli register -f ../provision/machines.yml
```

Notes:
- localhost rewrite: by default, `localhost`/`127.0.0.1` are rewritten to `host.docker.internal` when registering machines. This ensures the API (running in Docker) can reach host-published ports like `22221`. Disable with `--rewrite-localhost=false`.
- macOS: `host.docker.internal` works out of the box.
- Linux: your API/Scheduler containers must include:
  ```yaml
  extra_hosts:
    - "host.docker.internal:host-gateway"
  ```
  The CLI runs on the host and doesn’t need that mapping.

## Config (flags or env)

- `--api-base` (`API_BASE`, default `https://localhost`)
- `--api-prefix` (`API_PREFIX`, default `/api`)
- `--verify-tls` (`VERIFY_TLS`, default `false`)
- `--token-file` (`TOKEN_FILE`, default `~/.projet-iac/token.json`)
- `--rewrite-localhost` (`REWRITE_LOCALHOST`, default `true`)
- `--docker-host` (`DOCKER_HOST_GATEWAY_NAME`, default `host.docker.internal`)
- `--admin-user` (`ADMIN_DEFAULT_USERNAME`, default `admin`)
- `--admin-pass` (`ADMIN_DEFAULT_PASSWORD`, default `change-me`)
