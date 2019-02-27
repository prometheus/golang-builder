FROM        quay.io/prometheus/golang-builder:1.12-base
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN \
    apt-get update && apt-get install -y --no-install-recommends \
        crossbuild-essential-powerpc \
        crossbuild-essential-ppc64el \
    && rm -rf /var/lib/apt/lists/*

COPY rootfs /
