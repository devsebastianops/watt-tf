# Installation

We offer multiple ways to install Watt TF, depending on your preferences and environment.

## Installer
```bash
curl -sSL https://raw.githubusercontent.com/devsebastianops/watt-tf/main/install.sh | bash
```

## Via Go
```bash
go install github.com/devsebastianops/watt-tf/cmd/wtf@latest
```

## From Source
```bash
git clone https://github.com/devsebastianops/watt-tf.git
cd watt-tf
go build -o wtf ./cmd/wtf/main.go
```

## Docker
```bash
docker run --rm -v "$(pwd):/app" -w /app devsebastianops/watt-tf help
```

## Pre-built binaries
You can download the latest release from the [GitHub releases page](https://github.com/devsebastianops/watt-tf/releases).

## Verify your installation

```bash
wtf help
```