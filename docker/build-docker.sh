#!/bin/bash

set -e -x

VERSION=`cat package.json| grep version | awk -F':' '{print $2}'| tr -d "\", "`
COMMIT=`git rev-parse --short HEAD`

if [ ! -f dist/vcentercollector-${VERSION}-${COMMIT}.tar.gz ]
then
    echo "building binary...."
    npm run build:static
    go run build.go pkg-min-tar
else
    echo "skiping build..."
fi

export VERSION
export COMMIT

cp dist/vcentercollector-${VERSION}-${COMMIT}_${GOOS:-linux}_${GOARCH:-amd64}.tar.gz  docker/vcentercollector-last.tar.gz
cp conf/sample.vcentercollector.toml docker/vcentercollector.toml

cd docker

sudo docker build --label version="${VERSION}" --label commitid="${COMMIT}" -t tonimoreno/vcentercollector:${VERSION} -t tonimoreno/vcentercollector:latest .
rm vcentercollector-last.tar.gz
rm vcentercollector.toml

sudo docker push tonimoreno/vcentercollector:${VERSION}
sudo docker push tonimoreno/vcentercollector:latest
