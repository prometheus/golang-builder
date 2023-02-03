FROM goreleaser/goreleaser:v1.13.1 as goreleaser

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
        && git config --global --add safe.directory "*" \
        && curl -s -f -L https://deb.nodesource.com/gpgkey/nodesource.gpg.key -o /etc/apt/trusted.gpg.d/nodesource.asc \
        && echo "deb https://deb.nodesource.com/node_16.x/ bullseye main" > /etc/apt/sources.list.d/nodesource.list \
        && apt-get update \
        && apt-get install -y --no-install-recommends nodejs \
        && curl -s -f -L https://github.com/mikefarah/yq/releases/download/v4.30.5/yq_linux_amd64 -o "/bin/yq" \
        && echo "813d7d5b16a627a1fbbd763cb508fc8f143122f8f7195c872e7367d6295a331a /bin/yq" > /tmp/yq.sum \
        && sha256sum -c /tmp/yq.sum \
        && chmod  0755 /bin/yq \
        && rm -rf /tmp/yq.sum /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.19.5
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 36519702ae2fd573c9869461990ae550c8c0d955cd28d2827a6b159fda81ff95

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

RUN curl -s -f -L https://github.com/gotestyourself/gotestsum/releases/download/v1.7.0/gotestsum_1.7.0_linux_amd64.tar.gz -o gotestsum.tar.gz \
    && echo "b5c98cc408c75e76a097354d9487dca114996e821b3af29a0442aa6c9159bd40 gotestsum.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf gotestsum.tar.gz gotestsum \
    && rm gotestsum.tar.gz

VOLUME      /app
WORKDIR     /app
ENTRYPOINT  ["/builder.sh"]
