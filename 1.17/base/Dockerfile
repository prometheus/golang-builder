FROM goreleaser/goreleaser:v0.169.0 as goreleaser

FROM        debian:buster
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN \
    apt-get update \
        && apt-get dist-upgrade -y \
        && apt-get install -y --no-install-recommends \
            build-essential \
            bzr \
            ca-certificates \
            curl \
            git \
            gnupg \
            libsnmp-dev \
            make \
            unzip \
            yamllint \
        && curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - \
        && echo "deb https://dl.yarnpkg.com/debian/ stable main" > /etc/apt/sources.list.d/yarn.list \
        && curl -sS https://deb.nodesource.com/gpgkey/nodesource.gpg.key | apt-key add - \
        && echo "deb https://deb.nodesource.com/node_14.x/ buster main" > /etc/apt/sources.list.d/nodesource.list \
        && apt-get update \
        && apt-get install -y --no-install-recommends nodejs yarn \
        && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.17
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 6bf89fc4f5ad763871cf7eac80a2d594492de7a818303283f1366a7f6a30372d

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
    && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf golang.tar.gz \
    && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

COPY rootfs /

COPY --from=goreleaser /usr/local/bin/goreleaser /usr/local/bin/goreleaser

VOLUME      /app
WORKDIR     /app
ENTRYPOINT  ["/builder.sh"]
