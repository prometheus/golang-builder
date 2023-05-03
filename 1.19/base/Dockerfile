FROM goreleaser/goreleaser:v1.16.0 as goreleaser

FROM        debian:bullseye
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN \
    apt-get update \
        && apt-get full-upgrade -y \
        && apt-get install -y --no-install-recommends \
            build-essential \
            bzr \
            ca-certificates \
            curl \
            git \
            gnupg \
            jq \
            libsnmp-dev \
            make \
            unzip \
            yamllint \
            openssh-client \
        && git config --system --add safe.directory "*" \
        && curl -s -f -L https://deb.nodesource.com/gpgkey/nodesource.gpg.key -o /etc/apt/trusted.gpg.d/nodesource.asc \
        && echo "deb https://deb.nodesource.com/node_16.x/ bullseye main" > /etc/apt/sources.list.d/nodesource.list \
        && apt-get update \
        && apt-get install -y --no-install-recommends nodejs \
        && curl -s -f -L https://github.com/mikefarah/yq/releases/download/v4.31.2/yq_linux_amd64 -o "/bin/yq" \
        && echo "71ef4141dbd9aec3f7fb45963b92460568d044245c945a7390831a5a470623f7 /bin/yq" > /tmp/yq.sum \
        && sha256sum -c /tmp/yq.sum \
        && chmod  0755 /bin/yq \
        && rm -rf /tmp/yq.sum /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.19.9
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 e858173b489ec1ddbe2374894f52f53e748feed09dde61be5b4b4ba2d73ef34b

RUN curl -s -f -L "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
    && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf golang.tar.gz \
    && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

COPY rootfs /

COPY --from=goreleaser /usr/bin/goreleaser /usr/bin/goreleaser

RUN curl -s -f -L https://github.com/gotestyourself/gotestsum/releases/download/v1.9.0/gotestsum_1.9.0_linux_amd64.tar.gz -o gotestsum.tar.gz \
    && echo "0a0a310d6331447559b836407bcb4b98e758b9d4d3d829f3505f55ec74b26cae gotestsum.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf gotestsum.tar.gz gotestsum \
    && rm gotestsum.tar.gz

VOLUME      /app
WORKDIR     /app
ENTRYPOINT  ["/builder.sh"]
