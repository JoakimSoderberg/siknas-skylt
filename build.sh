#!/bin/bash
set -e

# TODO: Add command line options.

echo "### Building Docker image for Siknas-skylt Server (Backend) ###"
docker build -t siknas-skylt-server -f ./src/server/Dockerfile ./src/server

echo "### Building Dockeri mage for Aurelia web project (Frontend) ###"
docker build -t siknas-skylt-aurelia -f ./src/server/static/siknas-skylt/Dockerfile ./src/server/static/siknas-skylt

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

    docker run -it --rm -v $(pwd)/build/:/shared siknas-skylt-server \
        /bin/sh -c "GOOS=${os_arch[0]} GOARCH=${os_arch[1]} go build -v -o /shared/${output_dir}/siknas-skylt-server"
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
# TODO: Build debian package
# TODO: Systemd service files for Rpi
# TODO: Dependency on Xvfb and running that using Systemd
# TODO: Default settings file with static-path set to /usr/share ... by default.
