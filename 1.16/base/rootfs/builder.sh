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

  if [[ "${arch}" =~ ^armv.*$ ]]; then
    goarm=${arch##*v}
    arch="arm"

    echo "# ${goos}-${arch}v${goarm}"
    prefix=".build/${goos}-${arch}v${goarm}"
    mkdir -p "${prefix}"
    GOARM=${goarm} GOOS=${goos} GOARCH=${arch} make PREFIX="${prefix}" build
  else
    echo "# ${goos}-${arch}"
    prefix=".build/${goos}-${arch}"
    mkdir -p "${prefix}"
    GOOS=${goos} GOARCH=${arch} make PREFIX="${prefix}" build
  fi
done

exit 0
