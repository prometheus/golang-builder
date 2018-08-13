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
goarchs=(${goarchs[@]:-linux\/mips})
for goarch in "${goarchs[@]}"
do
  goos=${goarch%%/*}
  arch=${goarch##*/}

  echo "# ${goos}-${arch}"
  prefix=".build/${goos}-${arch}"
  mkdir -p "${prefix}"

  case "${arch}" in
    mips) cc="mips-linux-gnu-gcc" cxx="mips-linux-gnu-g++" ;;
    mips64) cc="mips64-linux-gnuabi64-gcc" cxx="mips64-linux-gnuabi64-g++" ;;
    mipsle) cc="mipsel-linux-gnu-gcc" cxx="mipsel-linux-gnu-g++" ;;
    mips64le) cc="mips64el-linux-gnuabi64-gcc" cxx="mips64el-linux-gnuabi64-g++" ;;
    *)
      echo 'Error: This is mips64/mips64le builder only.'
      exit 1
      ;;
  esac
  CC="${cc}" CXX="${cxx}" GOOS="${goos}" GOARCH="${arch}" \
    make PREFIX="${prefix}" build
done

exit 0
