#!/usr/bin/env bash
# Builds the release binaries of the game into dist/, one per system:
# battleship-<version>-<os>-<arch>, with .exe on Windows.
#
# Usage: scripts/build-release.sh <version>
#
# The version only names the files. The version the game shows with
# --versao comes from the git tag of the checkout, stamped by the Go tool.
set -euo pipefail

if [ "$#" -ne 1 ] || [ -z "$1" ]; then
  echo "uso: $0 <versão>" >&2
  exit 2
fi
version="$1"

cd "$(dirname "$0")/.."
rm -rf dist
mkdir -p dist

targets="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64"
for target in $targets; do
  goos="${target%/*}"
  goarch="${target#*/}"
  out="dist/battleship-${version}-${goos}-${goarch}"
  if [ "$goos" = "windows" ]; then
    out="${out}.exe"
  fi
  echo "→ ${out}"
  # Build the package (not main.go) so the Go tool stamps the version.
  CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
    go build -trimpath -ldflags "-s -w" -o "$out" .
done
