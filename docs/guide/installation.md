# Installation

We offer multiple ways to install Watt TF, depending on your preferences and environment.

## Via Go
```bash
go install github.com/devsebastianops/watt-tf/cmd/wtf@latest
wtf --help
```

## From Source
```bash
git clone https://github.com/devsebastianops/watt-tf.git
cd watt-tf
go build -o wtf ./cmd/wtf/main.go
./wtf --help
```

## Docker
```bash
docker run --rm -v "$(pwd):/app" -w /app devsebastianops/watt-tf --help
```

## Pre-built binaries
You can download the latest release from the [GitHub releases page](https://github.com/devsebastianops/watt-tf/releases).

