#!/usr/bin/env bash

# Copyright 2016 The Prometheus Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

source /common.sh

# Building binaries for the specified platforms
# The `build` Makefile target is required
declare -a goarchs
goarchs=(${goarchs[@]:-linux\/amd64})
for goarch in "${goarchs[@]}"
do
  goos=${goarch%%/*}
  arch=${goarch##*/}
  goarm='' cc='gcc' cxx='g++' extra_args=''

  case "${goos}" in
    windows)
      case "${arch}" in
        386) cc='i686-w64-mingw32-gcc' cxx='i686-w64-mingw32-g++' ;;
        amd64) cc='x86_64-w64-mingw32-gcc' cxx='x86_64-w64-mingw32-g++' ;;
      esac
      ;;
    darwin)
      case "${arch}" in
        amd64) cc='o64-clang' cxx='o64-clang++' extra_args='LD_LIBRARY_PATH=/usr/osxcross/lib' ;;
        arm64) cc='oa64-clang' cxx='oa64-clang++' extra_args='LD_LIBRARY_PATH=/usr/osxcross/lib' ;;
      esac
      ;;
    linux)
      case "${arch}" in
        386) cc='i686-linux-gnu-gcc' cxx='i686-linux-gnu-g++' ;;
        arm64) cc='aarch64-linux-gnu-gcc' cxx='aarch64-linux-gnu-g++' ;;
        armv*) goarm="${arch##*v}"
          if [[ "${goarm}" == "7" ]]; then
            cc='arm-linux-gnueabihf-gcc' cxx='arm-linux-gnueabihf-g++'
          else
            cc='arm-linux-gnueabi-gcc' cxx='arm-linux-gnueabi-g++'
          fi
          ;;
        mips) cc='mips-linux-gnu-gcc' cxx='mips-linux-gnu-g++' ;;
        mips64) cc='mips64-linux-gnuabi64-gcc' cxx='mips64-linux-gnuabi64-g++' ;;
        mipsle) cc='mipsel-linux-gnu-gcc' cxx='mipsel-linux-gnu-g++' ;;
        mips64le) cc='mips64el-linux-gnuabi64-gcc' cxx='mips64el-linux-gnuabi64-g++' ;;
        ppc64) cc='powerpc-linux-gnu-gcc' cxx='powerpc-linux-gnu-g++' ;;
        ppc64le) cc='powerpc64le-linux-gnu-gcc' cxx='powerpc64le-linux-gnu-g++' ;;
        s390x) cc='gcc-s390x-linux-gnu' cxx='g++-s390x-linux-gnu' ;;
      esac
      ;;
    netbsd)
      case "${arch}" in
        386) cc='i686-linux-gnu-gcc' cxx='i686-linux-gnu-g++' ;;
      esac
      ;;
  esac

  echo "# ${goos}-${arch}"
  prefix=".build/${goos}-${arch}"
  mkdir -p "${prefix}"

  echo "# Building with: CC='${cc}' CXX='${cxx}'"
  if [[ -n "${goarm}" ]]; then
    CC="${cc}" CXX="${cxx}" GOOS="${goos}" GOARCH="${arch}" GOARM="${goarm}" \
      make PREFIX="${prefix}" build ${extra_args}
  else
    CC="${cc}" CXX="${cxx}" GOOS="${goos}" GOARCH="${arch}" \
      make PREFIX="${prefix}" build ${extra_args}
  fi
done

exit 0
