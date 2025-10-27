# projet-iac-cli

Go CLI for the Projet IAC API at `https://localhost/api` (by default).

## Install (private repo)

Because this repository is private, the Go toolchain needs to fetch it directly from Git using your GitHub credentials. You have three easy options.

### Option A: go install (recommended for Go users)

1) Authenticate Git for GitHub access
- Easiest: use GitHub CLI and set up Git credentials:
  ```bash
  gh auth login
  gh auth setup-git
  ```
  Or use SSH keys (recommended) and verify:
  ```bash
  ssh -T git@github.com
  # Should say: "Hi <user>! You've successfully authenticated..."
  ```

2) Tell Go that the repo is private so it skips the public proxy and checksum DB:
```bash
go env -w GOPRIVATE=github.com/Jeomhps
```

3) Install the CLI:
- Latest:
  ```bash
  go install github.com/Jeomhps/projet-iac-cli@latest
  ```
- Specific version tag:
  ```bash
  go install github.com/Jeomhps/projet-iac-cli@v0.1.0
  ```

4) Make sure your Go bin is on PATH, then verify:
```bash
# macOS/Linux:
echo "$(go env GOPATH)/bin" | grep -q "$PATH" || echo "Add $(go env GOPATH)/bin to PATH"
projet-iac-cli --version
```

Notes
- If you prefer HTTPS with a Personal Access Token (PAT), run:
  ```bash
  gh auth login
  gh auth setup-git
  ```
  This configures Git to use your token automatically for github.com.
- The GOPRIVATE pattern can be scoped broadly (owner-wide) or to a single repo. Owner-wide is convenient:
  - `go env -w GOPRIVATE=github.com/Jeomhps`

### Option B: Download a release (no Go toolchain required)

1) Go to the project’s Releases page (you must be logged into GitHub):
- https://github.com/Jeomhps/projet-iac-cli/releases

2) Download the archive for your OS/CPU, then extract and place the binary on your PATH:
- macOS/Linux:
  ```bash
  tar -xzf projet-iac-cli_<version>_<os>_<arch>.tar.gz
  sudo mv projet-iac-cli /usr/local/bin/
  projet-iac-cli --version
  ```
- Windows:
  - Download the `.zip`, extract `projet-iac-cli.exe`, put it somewhere on your PATH.

### Option C: Build from source

- Using SSH (recommended):
  ```bash
  git clone git@github.com:Jeomhps/projet-iac-cli.git
  cd projet-iac-cli
  go build -o projet-iac-cli
  ```
- Using GitHub CLI:
  ```bash
  gh repo clone Jeomhps/projet-iac-cli
  cd projet-iac-cli
  go build -o projet-iac-cli
  ```

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
