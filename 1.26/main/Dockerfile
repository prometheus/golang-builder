FROM        quay.io/prometheus/golang-builder:1.26-base
LABEL       maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

ENV ZIG_VERSION=0.14.0
ENV ZIG_SHA256=473ec26806133cf4d1918caf1a410f8403a13d979726a9045b421b685031a982
ENV MACOS_SDK_PATH=/usr/local/lib/macos-sdk

RUN \
    apt-get update && apt-get install -y --no-install-recommends \
        mingw-w64 \
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
    && mkdir -p /usr/local/lib/zig \
    && tar -C /usr/local/lib/zig --strip-components=1 -xJf \
       "$(/bin/download.sh "https://ziglang.org/download/${ZIG_VERSION}/zig-linux-x86_64-${ZIG_VERSION}.tar.xz" \
          "${ZIG_SHA256}")" \
    && ln -s /usr/local/lib/zig/zig /usr/local/bin/zig

COPY macos-sdk-bin/ /usr/local/lib/macos-sdk/bin/
COPY macos-sdk/ /usr/local/lib/macos-sdk/
RUN chmod +x /usr/local/lib/macos-sdk/bin/*

WORKDIR /app

COPY rootfs /
