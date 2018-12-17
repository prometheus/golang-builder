FROM        debian:stable
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN \
    apt-get update \
        && apt-get dist-upgrade -y \
        && apt-get install -y --no-install-recommends \
            build-essential \
            ca-certificates \
            curl \
            git \
            gnupg \
            libsnmp-dev \
            make \
        && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.10.7
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 acbdedf28b55b38d2db6f06209a25a869a36d31bdcf09fd2ec3d40e1279e0592

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

COPY rootfs /

VOLUME      /app
WORKDIR     /app
ENTRYPOINT  ["/builder.sh"]
