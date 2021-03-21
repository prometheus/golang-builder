FROM        quay.io/prometheus/golang-builder:1.15-base
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN \
    apt-get update && apt-get install -y --no-install-recommends \
        clang \
        cmake \
        g++ \
        gcc \
        gcc-multilib \
        libc6-dev \
        libxml2-dev \
        lzma-dev \
        mingw-w64 \
        patch \
        xz-utils \
    && rm -rf /var/lib/apt/lists/*

ARG PROM_OSX_SDK_URL
ENV OSXCROSS_PATH=/usr/osxcross \
    OSXCROSS_REV=da2c3d4ff604458a931b08b3af800c5a454136de \
    SDK_VERSION=11.1 \
    DARWIN_VERSION=20.1 \
    OSX_VERSION_MIN=10.9
RUN \
    mkdir -p /tmp/osxcross && cd /tmp/osxcross \
    && curl -sSL "https://codeload.github.com/tpoechtrager/osxcross/tar.gz/${OSXCROSS_REV}" \
        | tar -C /tmp/osxcross --strip=1 -xzf - \
    && curl -sSLo tarballs/MacOSX${SDK_VERSION}.sdk.tar.xz "${PROM_OSX_SDK_URL}/MacOSX${SDK_VERSION}.sdk.tar.xz" \
    && UNATTENDED=yes ./build.sh >/dev/null \
    && mv target "${OSXCROSS_PATH}" \
    && rm -rf /tmp/osxcross "/usr/osxcross/SDK/MacOSX${SDK_VERSION}.sdk/usr/share/man"

ENV PATH $OSXCROSS_PATH/bin:$PATH

COPY rootfs /
