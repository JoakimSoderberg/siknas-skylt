#!/bin/bash

docker build -t siknas-skylt-debpack -f Dockerfile.debpack .

docker run -it --rm -v $(pwd):/shared -w /shared siknas-skylt-debpack \
    /bin/bash -c "VERBOSE=* go-bin-deb generate -a amd64 --version 0.0.1 -w /tmp/pkg-build/amd64 -o build/siknas-skylt-server.deb"
