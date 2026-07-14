FROM golang:1.25-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o bin/wtf ./cmd/wtf/main.go

FROM alpine:latest
COPY --from=build /app/bin/wtf /usr/bin/wtf
ENTRYPOINT ["/usr/bin/wtf"]