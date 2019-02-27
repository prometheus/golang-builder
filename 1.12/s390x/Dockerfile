FROM        quay.io/prometheus/golang-builder:1.12-base
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN \
    apt-get update && apt-get install -y --no-install-recommends \
        dpkg-cross \
        g++-s390x-linux-gnu \
        gcc-s390x-linux-gnu \
    && rm -rf /var/lib/apt/lists/*

COPY rootfs /
