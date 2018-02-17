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
goarchs=(${goarchs[@]:-linux\/ppc64})
for goarch in "${goarchs[@]}"
do
  goos=${goarch%%/*}
  arch=${goarch##*/}

  echo "# ${goos}-${arch}"
  prefix=".build/${goos}-${arch}"
  mkdir -p "${prefix}"

  if [[ "${arch}" == "ppc64" ]]; then
    CC="powerpc-linux-gnu-gcc" CXX="powerpc-linux-gnu-g++" GOOS=${goos} GOARCH=${arch} make PREFIX="${prefix}" build
  elif [[ "${arch}" == "ppc64le" ]]; then
    CC="powerpc64le-linux-gnu-gcc" CXX="powerpc64le-linux-gnu-g++" GOOS=${goos} GOARCH=${arch} make PREFIX="${prefix}" build
  else
    echo 'Error: This is mips/mipsel builder only.'
    exit 1
  fi
done

exit 0
