# Prometheus Golang builder Docker images

[![CircleCI](https://circleci.com/gh/prometheus/golang-builder/tree/master.svg?style=shield)][circleci]
[![Docker Stars](https://img.shields.io/docker/stars/prom/golang-builder.svg)][hub]
[![Docker Pulls](https://img.shields.io/docker/pulls/prom/golang-builder.svg)][hub]
[![Image Size](https://img.shields.io/imagelayers/image-size/prom/golang-builder/latest.svg)][imagelayers]
[![Image Layers](https://img.shields.io/imagelayers/layers/prom/golang-builder/latest.svg)][imagelayers]

## Details

Docker Builder Image for cross-building Golang Prometheus projects.

## Usage

Change the repository import path (`-i`) and target platforms (`-p`) according to your needs.
You can also use those images to run your tests by using the `-T` option.

```
Usage: builder.sh [args]
  -i,--import-path arg  : Go import path of the project
  -p,--platforms arg    : List of platforms (GOOS/GOARCH) to build separated by a space
  -T,--tests            : Go run tests then exit
```

### Requirements

This building process is using make to build and run tests. 
Therefore a `Makefile` with `build` and `test` targets is needed into the root of your source files.

### main/latest tag

```
docker run --rm -ti -v $(pwd):/app prom/golang-builder:main \
    -i "github.com/prometheus/alertmanager" \
    -p "linux/amd64 linux/386 darwin/amd64 darwin/386 windows/amd64 windows/386 freebsd/amd64 freebsd/386 openbsd/amd64 openbsd/386 netbsd/amd64 netbsd/386 dragonfly/amd64"
```

### arm tag

```
docker run --rm -ti -v $(pwd):/app prom/golang-builder:arm \
    -i "github.com/prometheus/alertmanager" \
    -p "linux/arm linux/arm64 freebsd/arm openbsd/arm netbsd/arm"
```

### powerpc tag

```
docker run --rm -ti -v $(pwd):/app prom/golang-builder:powerpc \
    -i "github.com/prometheus/alertmanager" \
    -p "linux/ppc64 linux/ppc64le"
```

### mips tag

mips/mipsel cross-build is currently not working.

```
docker run --rm -ti -v $(pwd):/app prom/golang-builder:mips \
    -i "github.com/prometheus/alertmanager" \
    -p "linux/mips linux/mipsel"
```

## Legal note

OSX/Darwin/Apple builds:
**[Please ensure you have read and understood the Xcode license
   terms before continuing.](https://www.apple.com/legal/sla/docs/xcode.pdf)**

## More information

  * You will find a Circle CI configuration in `circle.yml`.
  * All of the core developers are accessible via the [Prometheus Developers Mailinglist](https://groups.google.com/forum/?fromgroups#!forum/prometheus-developers) and the `#prometheus` channel on `irc.freenode.net`.

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md)

## License

Apache License 2.0, see [LICENSE](LICENSE).


[hub]: https://hub.docker.com/r/prom/golang-builder/
[circleci]: https://circleci.com/gh/prometheus/golang-builder
[imagelayers]: https://imagelayers.io/?images=prom/golang-builder:latest