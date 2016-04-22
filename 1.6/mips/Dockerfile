FROM        quay.io/prometheus/golang-builder:1.6-base
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

VOLUME  /app
WORKDIR /app

RUN \
    echo "deb http://emdebian.org/tools/debian/ unstable main" > /etc/apt/sources.list.d/crosstools.list \
    && curl http://emdebian.org/tools/debian/emdebian-toolchain-archive.key | apt-key add - \
    && apt-get update && apt-get install -y --no-install-recommends \
        crossbuild-essential-mipsel \
        g++-mips-linux-gnu \
        gcc-mips-linux-gnu \
    && rm -rf /var/lib/apt/lists/*

COPY rootfs /

ENTRYPOINT ["/builder.sh"]
