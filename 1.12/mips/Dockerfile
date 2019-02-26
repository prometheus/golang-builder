FROM        quay.io/prometheus/golang-builder:1.12-base
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN \
    apt-get update && apt-get install -y --no-install-recommends \
        crossbuild-essential-mipsel \
        g++-mips-linux-gnu \
        gcc-mips-linux-gnu \
        g++-mipsel-linux-gnu \
        gcc-mipsel-linux-gnu \
        g++-mips64-linux-gnuabi64 \
        gcc-mips64-linux-gnuabi64 \
        g++-mips64el-linux-gnuabi64 \
        gcc-mips64el-linux-gnuabi64 \
        libc6-dev-mipsel-cross \
        libc6-dev-mips64-cross \
        libc6-dev-mips64-mips-cross \
        libc6-dev-mips64-mipsel-cross \
        libc6-dev-mips64el-cross \
    && rm -rf /var/lib/apt/lists/*

COPY rootfs /
