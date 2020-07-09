#!/usr/bin/env bash
set -euo pipefail

if [[ "${LATEST:-false}" = "true" ]]; then
  gox -os="linux darwin windows" \
      -arch="386 amd64 arm" \
      -osarch='!darwin/arm' \
      -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" \
      -ldflags "-s -w -X main.Revision=$TRAVIS_COMMIT -X main.Version=${TRAVIS_TAG:-dev-build}" \
      -verbose \
      ./...

  pushd bin

  sha256sum s3fetch* | tee sha256sums
  export GZIP_OPT=-9
  for f in s3fetch* ; do
    zip -T9 "$f.zip" "$f"
    tar czvf "$f.tar.gz" --owner=0 --group=0 "$f"
    rm "$f"
  done

  popd
fi
