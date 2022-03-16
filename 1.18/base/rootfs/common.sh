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

set -eo pipefail

# This is a Makefile based building processus
[[ ! -e "./Makefile" ]] && echo "Error: A Makefile with 'build' and 'test' targets must be present into the root of your source files" && exit 1

usage() {
  base="$(basename "$0")"
  cat <<EOUSAGE
Usage: ${base} [args]
  -i,--import-path arg  : Go import path of the project
  -p,--platforms arg    : List of platforms (GOOS/GOARCH) to build separated by a space
  -T,--tests            : Go run tests then exit
EOUSAGE
}

if [[ $# -eq 0 ]]; then
  usage
fi

# Flag parsing
while [[ $# -gt 0 ]]; do
  opt="$1"
  case "${opt}" in
    -i|--import-path)
      repoName="$2"
      shift 2
      ;;
    -p|--platforms)
      IFS=' ' read -r -a goarchs <<< "$2"
      shift 2
      ;;
    -T|--tests)
      tests=1
      shift
      ;;
    *)
      echo "Error: Unknown option: ${opt}"
      usage
      exit 1
      ;;
  esac
done

[[ -z "${repoName}" ]] && echo "Error: {-i,--import-path} option is mandatory" && exit 1

# Get first path listed in GOPATH
goPath="${GOPATH%%:*}"
repoPath="${goPath}/src/${repoName}"

# Simulate the go src path with a symlink
mkdir -p "$(dirname "${repoPath}")"
ln -sf /app "${repoPath}"

# Running tests
# The `test` Makefile target is required
tests=${tests:-0}
if [[ ${tests} -eq 1 ]]; then
  # Need to be in the proper GOPATH to run tests
  cd "${repoPath}" ; make test
  exit 0
fi

# Look for the CGO envvar
CGO_ENABLED=${CGO_ENABLED:-0}
export CGO_ENABLED
