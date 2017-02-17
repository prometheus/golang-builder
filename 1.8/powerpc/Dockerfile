FROM        quay.io/prometheus/golang-builder:1.8-base
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN \
    echo "deb http://emdebian.org/tools/debian/ unstable main" > /etc/apt/sources.list.d/crosstools.list \
    && curl http://emdebian.org/tools/debian/emdebian-toolchain-archive.key | apt-key add - \
    && apt-get update && apt-get install -y --no-install-recommends \
        crossbuild-essential-powerpc \
        crossbuild-essential-ppc64el \
    && rm -rf /var/lib/apt/lists/*

COPY rootfs /
