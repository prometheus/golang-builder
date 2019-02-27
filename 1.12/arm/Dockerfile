FROM        quay.io/prometheus/golang-builder:1.12-base
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN \
    apt-get update && apt-get install -y --no-install-recommends \
        crossbuild-essential-arm64 \
        crossbuild-essential-armel \
        crossbuild-essential-armhf \
        linux-libc-dev-arm64-cross \
        linux-libc-dev-armel-cross \
        linux-libc-dev-armhf-cross \
    && rm -rf /var/lib/apt/lists/*

COPY rootfs /
