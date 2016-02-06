FROM        golang:1.5.3
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

VOLUME  /app
WORKDIR /app

ENV DEBIAN_FRONTEND=noninteractive
RUN \
    dpkg --add-architecture mips \
    && dpkg --add-architecture mipsel \
    && echo "deb http://emdebian.org/tools/debian/ jessie main" > /etc/apt/sources.list.d/crosstools.list \
    && curl http://emdebian.org/tools/debian/emdebian-toolchain-archive.key | apt-key add - \
    && apt-get update && apt-get update && apt-get upgrade -y \
    && apt-get install -y --no-install-recommends \
##        crossbuild-essential-mips \
        crossbuild-essential-mipsel \
    && rm -rf /var/lib/apt/lists/*

COPY rootfs /

ENTRYPOINT ["/builder.sh"]
