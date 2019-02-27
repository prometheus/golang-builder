# Prometheus Golang builder Docker images

[![CircleCI](https://circleci.com/gh/prometheus/golang-builder/tree/master.svg?style=shield)][circleci]
[![Docker Repository on Quay.io](https://quay.io/repository/prometheus/golang-builder/status)][quayio]

## Details

Docker Builder Image for cross-building Golang Prometheus projects.

- `latest`, `main`, `1.12-main`, `1.12.0-main` ([1.12/main/Dockerfile](1.12/main/Dockerfile))
- `arm`, `1.12-arm`, `1.12.0-arm` ([1.12/arm/Dockerfile](1.12/arm/Dockerfile))
- `powerpc`, `1.12-powerpc`, `1.12.0-powerpc` ([1.12/powerpc/Dockerfile](1.12/powerpc/Dockerfile))
- `mips`, `1.12-mips`, `1.12.0-mips` ([1.12/mips/Dockerfile](1.12/mips/Dockerfile))
- `s390x`, `1.12-s390x`, `1.12.0-s390x` ([1.12/s390x/Dockerfile](1.12/s390x/Dockerfile))
- `1.11-main`, `1.11.5-main` ([1.11/main/Dockerfile](1.11/main/Dockerfile))
- `arm`, `1.11-arm`, `1.11.5-arm` ([1.11/arm/Dockerfile](1.11/arm/Dockerfile))
- `powerpc`, `1.11-powerpc`, `1.11.5-powerpc` ([1.11/powerpc/Dockerfile](1.11/powerpc/Dockerfile))
- `mips`, `1.11-mips`, `1.11.5-mips` ([1.11/mips/Dockerfile](1.11/mips/Dockerfile))
- `s390x`, `1.11-s390x`, `1.11.5-s390x` ([1.11/s390x/Dockerfile](1.11/s390x/Dockerfile))

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
docker run --rm -ti -v $(pwd):/app quay.io/prometheus/golang-builder:main \
    -i "github.com/prometheus/prometheus" \
    -p "linux/amd64 linux/386 darwin/amd64 darwin/386 windows/amd64 windows/386 freebsd/amd64 freebsd/386 openbsd/amd64 openbsd/386 netbsd/amd64 netbsd/386 dragonfly/amd64"
```

### arm tag

```
docker run --rm -ti -v $(pwd):/app quay.io/prometheus/golang-builder:arm \
    -i "github.com/prometheus/prometheus" \
    -p "linux/arm linux/arm64 freebsd/arm openbsd/arm netbsd/arm"
```

### powerpc tag

```
docker run --rm -ti -v $(pwd):/app quay.io/prometheus/golang-builder:powerpc \
    -i "github.com/prometheus/prometheus" \
    -p "linux/ppc64 linux/ppc64le"
```

### mips tag

mips64/mips64le cross-build is currently available with golang 1.6.

```
docker run --rm -ti -v $(pwd):/app quay.io/prometheus/golang-builder:mips \
    -i "github.com/prometheus/prometheus" \
    -p "linux/mips64 linux/mips64le"
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

[quayio]: https://quay.io/repository/prometheus/golang-builder
[circleci]: https://circleci.com/gh/prometheus/golang-builder

