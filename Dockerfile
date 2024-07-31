# syntax=docker/dockerfile:1
FROM ghcr.io/jaredallard/scratch-cacerts:latest
ENTRYPOINT ["/usr/local/bin/dgcloud"]
COPY discord-gcloud /usr/local/bin/dgcloud
