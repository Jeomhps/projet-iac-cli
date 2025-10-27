# projet-iac-cli

Go CLI for the Projet IAC API at `https://localhost/api` (by default).

## Install (private via SSH)

Because this repository is private, the Go toolchain must fetch it via Git using your SSH key.

1) Ensure your SSH key works with GitHub
```bash
ssh -T git@github.com
# Expect: "Hi <you>! You've successfully authenticated..."
```

2) Tell Git to use SSH for any GitHub URL
```bash
git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"
# (Optional, narrower scope)
# git config --global url."ssh://git@github.com/Jeomhps/".insteadOf "https://github.com/Jeomhps/"
```

3) Tell Go that your repos are private (skip proxy/checksum DB)
```bash
go env -w GOPRIVATE=github.com/Jeomhps/*
```

4) Install with go
```bash
# Latest
go install github.com/Jeomhps/projet-iac-cli@latest

# Or a specific version tag
# go install github.com/Jeomhps/projet-iac-cli@v0.1.0
```

5) Ensure your Go bin is on PATH, then verify
```bash
# macOS/Linux:
echo "$(go env GOPATH)/bin" | grep -q "$PATH" || echo "Add $(go env GOPATH)/bin to PATH"
projet-iac-cli --version
```

Troubleshooting
- Remove/inspect the SSH mapping
  ```bash
  git config --global --get-all url."ssh://git@github.com/".insteadOf
  # To remove:
  # git config --global --unset-all url."ssh://git@github.com/".insteadOf
  ```
- If go install still prompts for HTTPS credentials, the mapping isn’t being applied; re-check step 2.

## Build

```bash
go mod tidy
go build -o projet-iac-cli
# optional: embed version/commit
# go build -ldflags "-X 'github.com/Jeomhps/projet-iac-cli/cmd.version=$(git describe --tags --always --dirty)' -X 'github.com/Jeomhps/projet-iac-cli/cmd.commit=$(git rev-parse --short HEAD)'" -o projet-iac-cli
```

## Quick start (dev)

```bash
export API_BASE=https://localhost
export API_PREFIX=/api
export VERIFY_TLS=false

./projet-iac-cli login          # prompts for username/password
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
- `--token-file` (`TOKEN_FILE`, default `~/.projet-iac/token.json`) — used if OS keychain is unavailable/disabled
- `--rewrite-localhost` (`REWRITE_LOCALHOST`, default `true`)
- `--docker-host` (`DOCKER_HOST_GATEWAY_NAME`, default `host.docker.internal`)
- `--keychain` (`KEYCHAIN`, default `auto`) — `auto|on|off` to control OS keychain use

See also: docs/KEYCHAIN.md for details on secure token storage.
