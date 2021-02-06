FROM alpine:latest
LABEL org.opencontainers.image.source=https://github.com/pbar1/atlantis-step-ghcodeowner

COPY bin/atlantis-step-ghcodeowner_linux_amd64 /usr/local/bin/atlantis-step-ghcodeowner

RUN chmod +x /usr/local/bin/atlantis-step-ghcodeowner

CMD ["/usr/local/bin/atlantis-step-ghcodeowner"]
