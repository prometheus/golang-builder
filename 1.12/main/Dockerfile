FROM        quay.io/prometheus/golang-builder:1.12-base
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN \
    dpkg --add-architecture i386 \
    && apt-get update && apt-get install -y --no-install-recommends \
        clang \
        g++ \
        gcc \
        gcc-multilib \
        libc6-dev \
        libc6-dev-i386 \
        linux-libc-dev:i386 \
        mingw-w64 \
        patch \
        xz-utils \
    && rm -rf /var/lib/apt/lists/*

ARG OSXCROSS_SDK_URL
ENV OSXCROSS_PATH=/usr/osxcross \
    OSXCROSS_REV=3034f7149716d815bc473d0a7b35d17e4cf175aa \
    SDK_VERSION=10.11 \
    DARWIN_VERSION=15 \
    OSX_VERSION_MIN=10.6
RUN \
    mkdir -p /tmp/osxcross && cd /tmp/osxcross \
    && curl -sSL "https://codeload.github.com/tpoechtrager/osxcross/tar.gz/${OSXCROSS_REV}" \
        | tar -C /tmp/osxcross --strip=1 -xzf - \
    && curl -sSLo tarballs/MacOSX${SDK_VERSION}.sdk.tar.xz ${OSXCROSS_SDK_URL} \
    && UNATTENDED=yes ./build.sh >/dev/null \
    && mv target "${OSXCROSS_PATH}" \
    && rm -rf /tmp/osxcross "/usr/osxcross/SDK/MacOSX${SDK_VERSION}.sdk/usr/share/man"

ENV PATH $OSXCROSS_PATH/bin:$PATH

COPY rootfs /
