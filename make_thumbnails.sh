#!/bin/bash

set -e

echo "NOTE! Assuming server is up and Xvfb is running inside of the server."
echo "(Also you need animations configured in the server config files, read the README.md)"
echo ""
echo "To start server:"
echo "  docker-compose up -d"
echo "  docker-compose exec server Xvfb :1 -screen 0, 1024x768x16 &  # Run in separate window"
echo ""
echo "Surf to port :8080 to make sure animations work before continuing!"
echo ""

read -p "Ready to continue? [y/n] " answer
if [ "${answer}" = "n" ]; then
    exit 0
fi

imagemagick_in_docker=0
read -p "Do you want to run ImageMagick in a docker? (slower but then no need to install it) [y/n] " answer
if [ "${answer}" = "y" ]; then
    imagemagick_in_docker=1
fi

pushd src/thumbgen/
docker build -t siknas-skylt-thumbgen .

if [ imagemagick_in_docker = 1 ]; then
    docker build -t siknas-skylt-gif2svg -f Dockerfile.gif2svg .
fi
popd

# Loop through and record all SVGs as animations.
docker run -it --rm \
    -v $(pwd)/src/server/static/siknas-skylt/images/animations:/go/src/app/output \
    siknas-skylt-thumbgen \
    sh -c "go run *.go --host $(docker-machine ip):8080 --output-frames --force"

# Make gifs from the SVGs
if [ imagemagick_in_docker = 1 ]; then
    docker run -it --rm \
        -v $(pwd)/src/server/static/siknas-skylt/images/animations:/go/src/app/output \
        siknas-skylt-gif2svg \
        ./makegifs.sh
else
    cd src/thumbgen && ./makegifs.sh -o ../server/static/siknas-skylt/images/animations
fi

