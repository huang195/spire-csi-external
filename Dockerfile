FROM golang:1-alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o spiffe-csi-driver ./cmd/spiffe-csi-driver

FROM alpine:latest as downloader
RUN apk add curl tar && curl -L -o /tmp/crictl.tar.gz \
  https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.31.1/crictl-v1.31.1-linux-amd64.tar.gz && \
  tar -C /tmp -xzf /tmp/crictl.tar.gz

FROM alpine:latest
COPY --from=builder /build/spiffe-csi-driver /bin/spiffe-csi-driver
COPY --from=downloader /tmp/crictl /bin/crictl
COPY --from=ghcr.io/spiffe/spire-agent:1.5.1 /opt/spire/bin/spire-agent /bin/spire-agent
RUN mkdir -p /var/run/secrets/spire

ENTRYPOINT ["/bin/spiffe-csi-driver"]
