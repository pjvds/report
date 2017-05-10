#!/bin/bash -e
if ! command -v git >/dev/null; then
    echo "cannot find git" >> /dev/stderr
    exit 1
fi

if ! command -v gox >/dev/null; then
    echo "cannot find gox, you can install it with:\n\n\rgo get github.com/mitchellh/gox" >> /dev/stderr
    exit 1
fi

VERSION=$(git describe --always || echo "unknown")
echo "building version: $VERSION"

build_dir=`mktemp -d`
function cleanup_build_dir {
  rm -rf "$build_dir"
}
trap cleanup_build_dir EXIT

gox -ldflags "-X main.version=$VERSION" \
    -osarch="linux/amd64" \
    -osarch="linux/386" \
    -osarch="windows/amd64" \
    -osarch="windows/386" \
    -osarch="darwin/amd64" \
    -osarch="darwin/386" \
    -output="$build_dir/slackme-$VERSION-{{.OS}}-{{.Arch}}/slackme" github.com/pjvds/slackme

if [ ! -d "./releases/" ]; then
  mkdir ./releases
fi

tar cfz ./releases/slackme-$VERSION-linux-amd64.tar.gz -C $build_dir/slackme-$VERSION-linux-amd64 .
tar cfz ./releases/slackme-$VERSION-linux-386.tar.gz -C $build_dir/slackme-$VERSION-linux-386 .
tar cfz ./releases/slackme-$VERSION-darwin-amd64.tar.gz -C $build_dir/slackme-$VERSION-darwin-amd64 .
tar cfz ./releases/slackme-$VERSION-darwin-386.tar.gz -C $build_dir/slackme-$VERSION-darwin-386 .
zip -qr ./releases/slackme-$VERSION-windows-amd64.zip $build_dir/slackme-$VERSION-windows-amd64
zip -qr ./releases/slackme-$VERSION-windows-386.zip $build_dir/slackme-$VERSION-windows-386

rm -rf $build_dir
