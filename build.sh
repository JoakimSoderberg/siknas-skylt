#!/bin/bash
set -e

# TODO: Add command line options.

echo "### Building Docker image for Siknas-skylt Server (Backend) ###"
docker build -t siknas-skylt-server -f ./src/server/Dockerfile ./src/server

echo "### Building Docker image for Aurelia web project (Frontend) ###"
docker build -t siknas-skylt-aurelia -f ./src/server/static/siknas-skylt/Dockerfile ./src/server/static/siknas-skylt

echo "### Building Docker image for thumbgen utility ###"
docker build -t siknas-skylt-thumbgen -f ./src/thumbgen/Dockerfile ./src/thumbgen

mkdir -p build

platforms=("windows-amd64" "linux-amd64" "linux-arm64" "linux-arm" "darwin-amd64")
for platform in "${platforms[@]}"
do
    echo "=== Building Siknas-skylt server for '${platform}' ==="

    os_arch=(${platform//-/ })
    mkdir -p ${platform}

    # Fix so go-bin-deb finds it, since the debian package names the architecture "armhf" instead of "arm".
    output_dir="${platform}"
    if [ "${platform}" == "linux-arm" ]; then
        output_dir="linux-armhf"
    fi

    # TODO: Support --quiet and remove -v here
    docker run -it --rm -v $(pwd)/build/:/shared siknas-skylt-server \
        /bin/sh -c "GOOS=${os_arch[0]} GOARCH=${os_arch[1]} go build -v -o /shared/${output_dir}/siknas-skylt-server"

    docker run -it --rm -v $(pwd)/build/:/shared siknas-skylt-thumbgen \
        /bin/sh -c "GOOS=${os_arch[0]} GOARCH=${os_arch[1]} go build -v -o /shared/${output_dir}/siknas-skylt-thumbgen"
done

mkdir -p build/static/

# Build static files for aurelia
docker run -it --rm \
    -v $(pwd)/src/server/static/siknas-skylt:/shared \
    -v $(pwd)/build/:/shared/build \
    -w /shared \
    siknas-skylt-aurelia \
    ./build.sh

# TODO: Build processing sketches
# TODO: Add travis

