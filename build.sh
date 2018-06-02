#!/bin/bash
set -e

docker build -t siknas-skylt-server -f ./src/server/Dockerfile ./src/server
docker build -t siknas-skylt-aurelia -f ./src/server/static/siknas-skylt/Dockerfile ./src/server/static/siknas-skylt

mkdir -p build

# TODO: Build multiple platforms
docker run -it --rm -v $(pwd)/build/:/shared siknas-skylt-server \
    /bin/sh -c "GOOS=windows GOARCH=amd64 go build -o /shared/siknas-skylt-server"

mkdir -p build/static/

# Build static files for aurelia
docker run -it --rm \
    -v $(pwd)/src/server/static/siknas-skylt:/shared \
    -v $(pwd)/build/:/shared/build \
    -w /shared \
    siknas-skylt-aurelia \
    ./build.sh



# TODO: Build processing sketches

# TODO: Build debian package
