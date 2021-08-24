#!/bin/bash

# Setting -x is absolutely forbidden as it could leak the GitHub token.
set -uo pipefail

GITHUB_MAIL="${GITHUB_MAIL:-prometheus-team@googlegroups.com}"
GITHUB_USER="${GITHUB_USER:-prombot}"
BRANCH_NAME="bump_version"

# GITHUB_TOKEN required scope: repo.repo_public
GITHUB_TOKEN="${GITHUB_TOKEN:-}"
if [[ -z ${GITHUB_TOKEN} ]]; then
  echo -e "\e[31mGitHub token (GITHUB_TOKEN) not set. Terminating.\e[0m"
  exit 1
fi

if ! GO111MODULE=on go run ./cmd/builder-bumper; then
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  # Check if a PR is already opened for the branch.
  prLink=$(curl --show-error --silent \
    -u "${GITHUB_USER}:${GITHUB_TOKEN}" \
    "https://api.github.com/repos/prometheus/golang-builder/pulls?head=prometheus:${BRANCH_NAME}" | jq '.[0].url')
  if [[ "${prLink}" != "null" ]]; then
    echo "Pull request already opened for ${BRANCH_NAME} branch: ${prLink}"
    echo "Either close it or merge it before running this script again!"
    exit 1
  fi

  git checkout -b "${BRANCH_NAME}"
  git config user.email "${GITHUB_MAIL}"
  git config user.name "${GITHUB_USER}"
  git add .
  git commit -s -m "Bump Go version"

  # Delete the remote branch in case it was merged but not deleted.
  # stdout and stderr are redirected to /dev/null otherwise git-push could leak the token in the logs.
  git push --quiet "https://${GITHUB_TOKEN}:@github.com/prometheus/golang-builder" ":${BRANCH_NAME}" 1>/dev/null 2>&1
  if git push --quiet "https://${GITHUB_TOKEN}:@github.com/prometheus/golang-builder" --set-upstream "${BRANCH_NAME}" 1>/dev/null 2>&1; then
    curl --show-error --silent \
      -u "${GITHUB_USER}:${GITHUB_TOKEN}" \
      -X POST \
      -d '{"title":"Bump Go version","base":"master","head":"'${BRANCH_NAME}'","body":""}' \
      "https://api.github.com/repos/prometheus/golang-builder/pulls"
  else
    echo "Failed to push changes to ${BRANCH_NAME}!"
    exit 1
  fi
fi
