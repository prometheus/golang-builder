FROM        quay.io/prometheus/golang-builder:1.8-base
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN \
    apt-get update && apt-get install -y --no-install-recommends \
        crossbuild-essential-mipsel \
        g++-mips-linux-gnu \
        gcc-mips-linux-gnu \
    && rm -rf /var/lib/apt/lists/*

COPY rootfs /
