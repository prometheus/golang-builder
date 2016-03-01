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

set -e

[ "$#" -lt 2 ] && echo "Missing args: $0 {VERSIONS} {VARIANTS}";

versions=( $1 )
variants=( $2 )

for version in "${versions[@]}"; do
  for variant in "${variants[@]}"; do
    (cd "${version}/${variant}"; make)
  done
done
