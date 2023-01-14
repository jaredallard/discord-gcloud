# syntax=docker/dockerfile:1.0-experimental
FROM golang:1.19.5 AS builder
ARG VERSION
ENV GOCACHE "/go-build-cache"
ENV CGO_ENABLED 0
WORKDIR /src

# Copy our source code into the container for building
COPY . .

# Cache dependencies across builds
RUN --mount=type=ssh --mount=type=cache,target=/go/pkg go mod download

# Build our application, caching the go build cache, but also using
# the dependency cache from earlier.
RUN --mount=type=ssh --mount=type=cache,target=/go/pkg --mount=type=cache,target=/go-build-cache \
  mkdir -p bin; \
  go build -o /src/bin/ -ldflags "-s -w"  -v ./cmd/...

FROM google/cloud-sdk:alpine
ENTRYPOINT ["/usr/bin/docker-entrypoint.sh"]

RUN apk add --no-cache bash 

COPY ./rootfs/ /
COPY --from=builder /src/bin/dgcloud /usr/local/bin/dgcloud