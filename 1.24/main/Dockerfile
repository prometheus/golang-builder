FROM rust:latest as rustbuilder

RUN cargo install --locked --git https://github.com/indygreg/apple-platform-rs --branch main --bin rcodesign apple-codesign

FROM        quay.io/prometheus/golang-builder:1.24-base
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN \
    apt-get update && apt-get install -y --no-install-recommends \
        clang \
        cmake \
        libc6-dev \
        libxml2-dev \
        lzma-dev \
        mingw-w64 \
        patch \
        xz-utils \
        crossbuild-essential-arm64 linux-libc-dev-arm64-cross \
        crossbuild-essential-armel linux-libc-dev-armel-cross \
        crossbuild-essential-armhf linux-libc-dev-armhf-cross \
        crossbuild-essential-i386 linux-libc-dev-i386-cross \
        crossbuild-essential-mips linux-libc-dev-mips-cross \
        crossbuild-essential-mipsel linux-libc-dev-mipsel-cross \
        crossbuild-essential-mips64 linux-libc-dev-mips64-cross \
        crossbuild-essential-mips64el linux-libc-dev-mips64el-cross \
        crossbuild-essential-powerpc linux-libc-dev-powerpc-cross \
        crossbuild-essential-ppc64el linux-libc-dev-ppc64el-cross \
        crossbuild-essential-s390x linux-libc-dev-s390x-cross \
        gcc-riscv64-linux-gnu g++-riscv64-linux-gnu libc6-dev-riscv64-cross \
        libc6-riscv64-cross linux-libc-dev-riscv64-cross \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /tmp/osxcross
    
COPY --from=rustbuilder /usr/local/cargo/bin/rcodesign /usr/local/bin/rcodesign

ARG PROM_OSX_SDK_URL
ENV OSXCROSS_PATH=/usr/osxcross \
    OSXCROSS_REV=ff8d100f3f026b4ffbe4ce96d8aac4ce06f1278b \
    SDK_VERSION=14 \
    DARWIN_VERSION=23 \
    OSX_VERSION_MIN=10.13

WORKDIR /tmp/osxcross
RUN \
    curl -s -f -L "https://codeload.github.com/tpoechtrager/osxcross/tar.gz/${OSXCROSS_REV}" \
      | tar -C /tmp/osxcross --strip=1 -xzf - \
    && curl -s -f -L "${PROM_OSX_SDK_URL}/MacOSX${SDK_VERSION}.sdk.tar.xz" -o "tarballs/MacOSX${SDK_VERSION}.sdk.tar.xz" \
    && UNATTENDED=yes JOBS=2 ./build.sh \
    && mv target "${OSXCROSS_PATH}" \
    && rm -rf /tmp/osxcross "/usr/osxcross/SDK/MacOSX${SDK_VERSION}.sdk/usr/share/man"

WORKDIR /app

ENV PATH $OSXCROSS_PATH/bin:$PATH

COPY rootfs /
