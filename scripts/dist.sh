#!/usr/bin/env bash
# Builds the release archives for every platform into dist/.
# Usage: scripts/dist.sh [version]
# OPENSUBTITLES_API_KEY, when set, is built into the binaries.
set -euo pipefail
cd "$(dirname "$0")/.."

version=${1:-0.0.0-dev}
ldflags="-s -w -X main.version=$version"
if [ -n "${OPENSUBTITLES_API_KEY:-}" ]; then
  ldflags="$ldflags -X main.defaultAPIKey=$OPENSUBTITLES_API_KEY"
else
  echo "warning: OPENSUBTITLES_API_KEY not set, users will need their own key" >&2
fi

rm -rf dist
mkdir -p dist/stage

build() { # goos goarch output [goarm]
  CGO_ENABLED=0 GOOS=$1 GOARCH=$2 GOARM=${4:-} go build -trimpath -ldflags "$ldflags" -o "$3" .
}

pack() { # platform binary
  local stage=dist/stage/$1 name=subtitlewatcher-$version-$1
  mkdir -p "$stage"
  mv "$2" "$stage/"
  cp README.md LICENSE "$stage/"
  case $1 in
    linux-* | freebsd-*) tar -czf "dist/$name.tar.gz" -C "$stage" . ;;
    *) (cd "$stage" && zip -qr "../../$name.zip" .) ;;
  esac
}

# macOS: one universal binary when lipo is available (macOS hosts).
build darwin arm64 dist/stage/mac-arm64
build darwin amd64 dist/stage/mac-amd64
if command -v lipo >/dev/null; then
  lipo -create -output dist/stage/subtitlewatcher dist/stage/mac-arm64 dist/stage/mac-amd64
  pack macos dist/stage/subtitlewatcher
else
  mv dist/stage/mac-arm64 dist/stage/subtitlewatcher && pack macos-arm64 dist/stage/subtitlewatcher
  mv dist/stage/mac-amd64 dist/stage/subtitlewatcher && pack macos-x64 dist/stage/subtitlewatcher
fi
rm -f dist/stage/mac-*

while read -r goos goarch platform goarm; do
  bin=dist/stage/subtitlewatcher
  [ "$goos" = windows ] && bin=$bin.exe
  build "$goos" "$goarch" "$bin" "$goarm"
  pack "$platform" "$bin"
done <<'TARGETS'
windows amd64 windows-x64
windows arm64 windows-arm64
linux amd64 linux-x64
linux arm64 linux-arm64
linux arm linux-armv7 7
freebsd amd64 freebsd-x64
TARGETS

rm -rf dist/stage
(cd dist && shasum -a 256 ./*.zip ./*.tar.gz | sed 's| \./| |' > checksums.txt)
ls -lh dist
