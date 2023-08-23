FROM goreleaser/goreleaser:v1.21.2 as goreleaser

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
        && curl -s -f -L https://github.com/mikefarah/yq/releases/download/v4.34.2/yq_linux_amd64 -o "/bin/yq" \
        && echo "1952f93323e871700325a70610d2b33bafae5fe68e6eb4aec0621214f39a4c1e /bin/yq" > /tmp/yq.sum \
        && sha256sum -c /tmp/yq.sum \
        && chmod  0755 /bin/yq \
        && rm -rf /tmp/yq.sum /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.21.0
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 d0398903a16ba2232b389fb31032ddf57cac34efda306a0eebac34f0965a0742

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

RUN curl -s -f -L https://github.com/gotestyourself/gotestsum/releases/download/v1.10.1/gotestsum_1.10.1_linux_amd64.tar.gz -o gotestsum.tar.gz \
    && echo "44be2c02d4cf99cdd61edcb27851ef98ef8724a2ae3355b438bd108e9abb9056 gotestsum.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf gotestsum.tar.gz gotestsum \
    && rm gotestsum.tar.gz

VOLUME      /app
WORKDIR     /app
ENTRYPOINT  ["/builder.sh"]
