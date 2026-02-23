FROM goreleaser/goreleaser:v2.13.3 as goreleaser

FROM        debian:bookworm
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

COPY download.sh /bin/download.sh

ENV YQ_VERSION 4.52.2
ENV YQ_URL https://github.com/mikefarah/yq/releases/download/v${YQ_VERSION}/yq_linux_amd64
ENV YQ_SUM a74bd266990339e0c48a2103534aef692abf99f19390d12c2b0ce6830385c459

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
        && curl -s -f -L https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key -o /etc/apt/nodesource.asc \
        && echo "deb [signed-by=/etc/apt/nodesource.asc] https://deb.nodesource.com/node_22.x nodistro main" > /etc/apt/sources.list.d/nodesource.list \
        && apt-get update \
        && apt-get install -y --no-install-recommends nodejs \
        && yq_file=$(/bin/download.sh ${YQ_URL} ${YQ_SUM}) \
        && mv -v ${yq_file} /bin/yq \
        && chmod  0755 /bin/yq \
        && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.26.0
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 aac1b08a0fb0c4e0a7c1555beb7b59180b05dfc5a3d62e40e9de90cd42f88235

RUN go_file=$(/bin/download.sh ${GOLANG_DOWNLOAD_URL} ${GOLANG_DOWNLOAD_SHA256}) \
    && tar -C /usr/local -xzf ${go_file} \
    && rm -v ${go_file}

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

COPY rootfs /

COPY --from=goreleaser /usr/bin/goreleaser /usr/bin/goreleaser

ENV GOTESTSUM_VERSION 1.13.0
ENV GOTESTSUM_URL https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_linux_amd64.tar.gz
ENV GOTESTSUM_SUM 11ccddeaf708ef228889f9fe2f68291a75b27013ddfc3b18156e094f5f40e8ee

RUN gotestsum_file=$(/bin/download.sh ${GOTESTSUM_URL} ${GOTESTSUM_SUM}) \
    && tar -C /usr/local -xzf ${gotestsum_file} gotestsum \
    && rm ${gotestsum_file}

VOLUME      /app
WORKDIR     /app
ENTRYPOINT  ["/builder.sh"]
