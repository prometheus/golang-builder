FROM        golang:1.5.3
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

VOLUME  /app
WORKDIR /app

ENV DEBIAN_FRONTEND=noninteractive
RUN \
    dpkg --add-architecture arm64 \
    && dpkg --add-architecture armel \
    && dpkg --add-architecture armhf \
    && echo "deb http://emdebian.org/tools/debian/ jessie main" > /etc/apt/sources.list.d/crosstools.list \
    && curl http://emdebian.org/tools/debian/emdebian-toolchain-archive.key | apt-key add - \
    && apt-get update && apt-get upgrade -y \
    && apt-get install -y --no-install-recommends \
        crossbuild-essential-arm64 \
        crossbuild-essential-armel \
        crossbuild-essential-armhf \
    && rm -rf /var/lib/apt/lists/*

COPY rootfs /

ENTRYPOINT ["/builder.sh"]
