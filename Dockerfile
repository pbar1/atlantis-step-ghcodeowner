FROM alpine:latest

COPY bin/atlantis-step-ghcodeowner_linux_amd64 /usr/local/bin/atlantis-step-ghcodeowner

RUN chmod +x /usr/local/bin/atlantis-step-ghcodeowner

CMD ["/usr/local/bin/atlantis-step-ghcodeowner"]
