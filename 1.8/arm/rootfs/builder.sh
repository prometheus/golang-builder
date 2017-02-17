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
goarchs=(${goarchs[@]:-linux\/arm})
for goarch in "${goarchs[@]}"
do
  goos=${goarch%%/*}
  arch=${goarch##*/}

  if [ "${arch}" = "arm" ]; then
    goarms=(5 6 7)
    for goarm in "${goarms[@]}"
    do
      echo "# ${goos}-${arch}v${goarm}"
      prefix=".build/${goos}-${arch}v${goarm}"
      mkdir -p "${prefix}"

      if [ "${goarm}" = 7 ]; then
        CC="arm-linux-gnueabihf-gcc" CXX="arm-linux-gnueabihf-g++" GOARM=${goarm} GOOS=${goos} GOARCH=${arch} make PREFIX="${prefix}" build
      else
        CC="arm-linux-gnueabi-gcc" CXX="arm-linux-gnueabi-g++" GOARM=${goarm} GOOS=${goos} GOARCH=${arch} make PREFIX="${prefix}" build
      fi
    done
  elif [ "${arch}" = "arm64" ]; then
    echo "# ${goos}-${arch}"
    prefix=".build/${goos}-${arch}"
    mkdir -p "${prefix}"

    CC="aarch64-linux-gnu-gcc" CXX="aarch64-linux-gnu-g++" GOOS=${goos} GOARCH=${arch} make PREFIX="${prefix}" build
  else
    echo 'Error: This is arm/arm64 builder only.'
  fi
done

exit 0
