#!/bin/bash

# Setting -x is absolutely forbidden as it could leak the GitHub token.
set -uo pipefail

GITHUB_MAIL="${GITHUB_MAIL:-prometheus-team@googlegroups.com}"
GITHUB_USER="${GITHUB_USER:-prombot}"
BRANCH_NAME="${BRANCH_NAME:-bump_version}"
DEFAULT_BRANCH="${DEFAULT_BRANCH:-$(git remote show origin | awk -F': ' '$1 == "  HEAD branch" {print $2}')}"

ORG="${ORG:-prometheus}"
REPO="${REPO:-golang-builder}"

# GITHUB_TOKEN required scope: repo.repo_public
GITHUB_TOKEN="${GITHUB_TOKEN:-}"
if [[ -z ${GITHUB_TOKEN} ]]; then
  echo -e "\e[31mGitHub token (GITHUB_TOKEN) not set. Terminating.\e[0m"
  exit 1
fi

if ! GO111MODULE=on go run ./cmd/builder-bumper; then
  exit 1
fi

## Internal functions
github_api() {
  local url
  url="https://api.github.com/${1}"
  shift 1
  curl --retry 5 --silent --fail -u "${GITHUB_USER}:${GITHUB_TOKEN}" "${url}" "$@"
}

if [[ -n "$(git status --porcelain)" ]]; then
  # Check if a PR is already opened for the branch.
  prLink=$(github_api "repos/${ORG}/${REPO}/pulls?state=open&head=${ORG}:${BRANCH_NAME}" | jq '.[0].html_url')
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
  git push --quiet "https://${GITHUB_TOKEN}:@github.com/${ORG}/${REPO}" ":${BRANCH_NAME}" 1>/dev/null 2>&1
  if git push --quiet "https://${GITHUB_TOKEN}:@github.com/${ORG}/${REPO}" --set-upstream "${BRANCH_NAME}" 1>/dev/null 2>&1; then
    post_json="$(printf '{"title":"Bump Go version","base":"%s","head":"%s","body":""}' "${DEFAULT_BRANCH}" "${BRANCH_NAME}")"
    github_api "repos/${ORG}/${REPO}/pulls" \
      -X POST \
      --data "${post_json}"
  else
    echo "Failed to push changes to ${BRANCH_NAME}!"
    exit 1
  fi
fi
