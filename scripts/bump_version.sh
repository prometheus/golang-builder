#!/bin/bash

# Setting -x is absolutely forbidden as it could leak the GitHub token.
set -uo pipefail

GITHUB_MAIL="${GITHUB_MAIL:-prometheus-team@googlegroups.com}"
GITHUB_USER="${GITHUB_USER:-prombot}"

# GITHUB_TOKEN required scope: repo.repo_public
GITHUB_TOKEN="${GITHUB_TOKEN:-}"
if [[ -z ${GITHUB_TOKEN} ]]; then
  echo -e "\e[31mGitHub token (GITHUB_TOKEN) not set. Terminating.\e[0m"
  exit 1
fi

if ! go run ./cmd/builder-bumper; then
  exit 1
fi

if [[ -n $(git status --porcelain) ]]; then
  git checkout -b bump_version
  git config user.email "${GITHUB_MAIL}"
  git config user.name "${GITHUB_USER}"
  git add .
  git commit -s -m "Bump Go version"
  # stdout and stderr are redirected to /dev/null otherwise git-push could leak the token in the logs.
  if git push --quiet "https://${GITHUB_TOKEN}:@github.com/prometheus/golang-builder" --set-upstream bump_version 1>/dev/null 2>&1; then
    curl --show-error --silent \
      -u "${GITHUB_USER}:${GITHUB_TOKEN}" \
      -X POST \
      -d "{\"title\":\"Bump Go version\",\"base\":\"master\",\"head\":\"bump_version\",\"body\":\"\"}" \
      "https://api.github.com/repos/prometheus/golang-builder/pulls"
  fi
fi
