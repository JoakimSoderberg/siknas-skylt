#!/bin/bash

docker build -t siknas-skylt-debpack -f Dockerfile.debpack .

platforms=("amd64-amd64" "arm-armhf")
for platform in "${platforms[@]}"
do
    goarch_debarch=(${platform//-/ })

    echo "=== Building Siknas-skylt debian package for '${goarch_debarch[1]}' ==="

    docker run -it --rm -v $(pwd):/shared -w /shared siknas-skylt-debpack \
    /bin/bash -c "VERBOSE=* go-bin-deb generate -a ${goarch_debarch[1]} --version 0.0.1 -w /tmp/pkg-build/${goarch_debarch[0]} -o build/siknas-skylt-server-${goarch_debarch[1]}.deb"
done
