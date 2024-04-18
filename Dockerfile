# CI system runs using amd64, so this allows us to not need QEMU, but
# will break builds on non-amd64 Linux systems (sorry).
FROM --platform=amd64 alpine:3.19 as cacerts
RUN apk add --no-cache ca-certificates

FROM alpine:3.19
ENTRYPOINT ["/usr/local/bin/dgcloud"]

COPY --from=cacerts /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY discord-gcloud /usr/local/bin/dgcloud
