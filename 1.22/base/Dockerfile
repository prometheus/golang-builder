FROM goreleaser/goreleaser:v1.20.0 as goreleaser

FROM        debian:bookworm
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

COPY download.sh /bin/download.sh

ENV YQ_VERSION 4.35.1
ENV YQ_URL https://github.com/mikefarah/yq/releases/download/v${YQ_VERSION}/yq_linux_amd64
ENV YQ_SUM bd695a6513f1196aeda17b174a15e9c351843fb1cef5f9be0af170f2dd744f08

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
        && echo "deb https://deb.nodesource.com/node_20.x/ bookworm main" > /etc/apt/sources.list.d/nodesource.list \
        && apt-get update \
        && apt-get install -y --no-install-recommends nodejs \
        && yq_file=$(/bin/download.sh ${YQ_URL} ${YQ_SUM}) \
        && mv -v ${yq_file} /bin/yq \
        && chmod  0755 /bin/yq \
        && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.22.2
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 5901c52b7a78002aeff14a21f93e0f064f74ce1360fce51c6ee68cd471216a17

RUN go_file=$(/bin/download.sh ${GOLANG_DOWNLOAD_URL} ${GOLANG_DOWNLOAD_SHA256}) \
    && tar -C /usr/local -xzf ${go_file} \
    && rm -v ${go_file}

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

COPY rootfs /

COPY --from=goreleaser /usr/bin/goreleaser /usr/bin/goreleaser

ENV GOTESTSUM_VERSION 1.7.0
ENV GOTESTSUM_URL https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_linux_amd64.tar.gz
ENV GOTESTSUM_SUM b5c98cc408c75e76a097354d9487dca114996e821b3af29a0442aa6c9159bd40

RUN gotestsum_file=$(/bin/download.sh ${GOTESTSUM_URL} ${GOTESTSUM_SUM}) \
    && tar -C /usr/local -xzf ${gotestsum_file} gotestsum \
    && rm ${gotestsum_file}

VOLUME      /app
WORKDIR     /app
ENTRYPOINT  ["/builder.sh"]
