# gwk

**G**it **W**al**k** — Walk multiple git repos and inspect their latest commits.

A read-only multi-repo inspector that fetches and pulls the latest from a set of hardcoded git URLs.

## Quick Start

```bash
# Install tools (mage)
go generate ./...

# Build for current platform
mage build

# Run
./bin/gwk

# Clean build artifacts
mage clean
```

## Mage Commands

- `mage build` — Run vet, then build to `bin/gwk` with version info baked in
- `mage install` — Build and copy binary to `~/.bio/bin/gwk`
- `mage clean` — Remove build artifacts
- `mage vet` — Run `go vet ./...`

## Usage

```
gwk inspect    Walk all configured repos
gwk version    Print version info
```
