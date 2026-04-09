# speedy-cli

A fast, fun, colorful terminal toolkit for:
- internet speed checks
- port conflict detection
- environment variable validation
- NEPSE stock updates
- mini ASCII animations and graphs

## Install

### Build locally

```bash
go mod tidy
go build -o speedy-cli .
./speedy-cli --help
```

### Homebrew (tap)

```bash
brew install samirkoirala/homebrew-tap/speedy-cli
```

### APT (.deb)

```bash
sudo apt install ./speedy-cli_<version>_linux_amd64.deb
```

Use `.goreleaser.yaml` and GitHub Actions release workflow to produce Homebrew formula updates and `.deb` packages.

## Open Source Standards

- License: `MIT` (`LICENSE`)
- Security disclosures: `SECURITY.md`
- Contribution process: `CONTRIBUTING.md`
- Community expectations: `CODE_OF_CONDUCT.md`
- Dependency updates: `.github/dependabot.yml`
- CI and vulnerability checks: `.github/workflows/ci.yml`

## Release (Homebrew + APT)

1. Push code to `main`
2. Create and push a semver tag
3. GitHub Actions runs GoReleaser and publishes release artifacts

```bash
git tag v0.1.0
git push origin main --tags
```

Expected release outputs:
- GitHub release with binaries and checksums
- Homebrew tap formula update (to `samirkoirala/homebrew-tap`)
- Debian package artifact (`.deb`) for apt-style installation

## Commands

### 1) Check internet speed

```bash
speedy-cli check-speed --graph
```

Example:

```text
🌐 Speed Test Results:
Download: 123 Mbps ███████████
Upload:   45 Mbps ██████
Ping:     12 ms ⚡
```

### 2) Check ports

```bash
speedy-cli ports --port 3000
```

### 3) Validate env vars

```bash
speedy-cli env --file .env
```

### 4) Fun mode

```bash
speedy-cli fun-mode
```

### 5) NEPSE stocks

```bash
speedy-cli stocks --top 5
speedy-cli stocks --symbol NABIL
```

## Global flags

- `--verbose` show extra details
- `--json` machine-readable output

## Project layout


```text
speedy-cli/
├── cmd/
│   ├── root.go
│   ├── speed.go
│   ├── ports.go
│   ├── env.go
│   ├── fun.go
│   └── stocks.go
├── internal/
│   ├── speedtest/
│   ├── portcheck/
│   ├── envcheck/
│   ├── animation/
│   └── stocks/
├── main.go
└── go.mod
```
