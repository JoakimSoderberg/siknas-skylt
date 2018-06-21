#!/bin/bash
set -e

show_help() {
    cat << EOF
Build Siknas-skylt server
-------------------------
Commands:

  docker-images
    Build all docker images needed.

  docker-{server,aurelia,thumbgen}
    Builds only the specific docker image.


  all (default)
    Build everything (but not the docker-images!)

  go
    Build the go code server/thumbgen

  static
    Builds the static files

  animations
    Builds the animations

  thumbs
    Build the thumbnails for the animations (SLOW)

  debian
    Builds the debian package (requires the other things to be built)
  help
    Show this help text and exit.
EOF
}

# TODO: Variable for build-dir

build_go_executables() {
    mkdir -p build

    # TODO: Enable building just specific build.

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
}

build_aurelia() {
    mkdir -p build/static/

    # Build static files for aurelia
    docker run -it --rm \
        -v $(pwd)/src/server/static/siknas-skylt:/shared \
        -v $(pwd)/build/:/shared/build \
        -w /shared \
        siknas-skylt-aurelia \
        ./build.sh
}

build_docker_server() {
    echo "### Building Docker image for Siknas-skylt Server (Backend) ###"
    docker build -t siknas-skylt-server -f ./src/server/Dockerfile ./src/server
}

build_docker_aurelia() {
    echo "### Building Docker image for Aurelia web project (Frontend) ###"
    docker build -t siknas-skylt-aurelia -f ./src/server/static/siknas-skylt/Dockerfile ./src/server/static/siknas-skylt
}

build_docker_thumbgen() {
    echo "### Building Docker image for thumbgen utility ###"
    docker build -t siknas-skylt-thumbgen -f ./src/thumbgen/Dockerfile ./src/thumbgen
}

build_thumbgen() {
    # TODO: Break this out nicer
    build_go_executables
}

build_debian() {
    # TODO: Make nicer.
    ./build-debian.sh
}

build_thumbs() {
    echo "Not supported yet"
    # TODO: Support building thumbnails
}

build_animations() {
    echo "Not supported yet"
    # TODO: Support building animations
}

# TODO: Add listing gifs and animations
case "$1" in
    "all")
        build_go_executables
        build_aurelia
        build_debian
        ;;
    "docker-images")
        shift
        build_docker_server
        build_docker_aurelia
        build_docker_thumbgen
        ;;
    "docker-server")
        shift
        build_docker_server $@
        ;;
    "docker-aurelia")
        shift
        build_docker_aurelia $@
        ;;
    "docker-thumbgen")
        shift
        build_docker_thumbgen $@
        ;;
    "go"|"app"|"server")
        shift
        build_go_executables $@
        ;;
    "aurelia"|"static")
        shift
        build_aurelia $@
        ;;
    "debian")
        shift
        build_debian $@
        ;;
    "animations")
        shift
        build_animations $@
        ;;
    "thumbs"|"thumbnails")
        shift
        build_thumbs $@
        ;;
    "thumbgen")
        shift
        build_thumbgen $@
        ;;
    ""|"help"|"-h"|"--help")
        show_help
        exit 0
        ;;
    *)
        show_help
        echo
        echo "Unknown command '$1'"
        exit 1
        ;;
esac

# TODO: Build processing sketches
# TODO: After building the program run 
# TODO: Add travis
