ARG BUILD_TIMESTAMP
ARG BUILD_COMMIT
ARG BUILD_VERSION

FROM golang:1.25-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o bin/wtf -ldflags "-X github.com/devsebastianops/watt-tf/internal/cli.Version=${BUILD_VERSION} \
    -X github.com/devsebastianops/watt-tf/internal/cli.Commit=${BUILD_COMMIT} \
    -X github.com/devsebastianops/watt-tf/internal/cli.BuildTime=${BUILD_TIMESTAMP}" \
    ./cmd/wtf/main.go

FROM alpine:latest
COPY --from=build /app/bin/wtf /usr/bin/wtf
ENTRYPOINT ["/usr/bin/wtf"]