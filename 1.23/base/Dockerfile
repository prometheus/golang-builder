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
        && curl -fsSL https://deb.nodesource.com/setup_22.x -o nodesource_setup.sh \
        && bash nodesource_setup.sh \
        && apt-get update \
        && apt-get install -y --no-install-recommends nodejs \
        && yq_file=$(/bin/download.sh ${YQ_URL} ${YQ_SUM}) \
        && mv -v ${yq_file} /bin/yq \
        && chmod  0755 /bin/yq \
        && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.23.6
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 9379441ea310de000f33a4dc767bd966e72ab2826270e038e78b2c53c2e7802d

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
