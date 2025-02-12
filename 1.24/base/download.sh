#!/usr/bin/env bash
#
# Description: Download a file and verify SHA256.

set -u -o pipefail

if [[ $# -ne 2 ]] ; then
  echo "usage: $(basename $0) <url> <checksum>"
  exit 1
fi

url="$1"
sum="$2"

outfile=$(mktemp)

if ! curl -s -f -L "${url}" -o "${outfile}" ; then
  echo "ERROR: Failed to download ${url} to ${outfile}" 1>&2
  rm "${outfile}"
  exit 1
fi

if ! echo "${sum} ${outfile}" | sha256sum -c - > /dev/null ; then
  echo "ERROR: Checksum failed" 1>&2
  rm "${outfile}"
  exit 1
fi

echo "${outfile}"
