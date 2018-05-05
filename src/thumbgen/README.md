Siknäs-skylt Thumbnail Generator
================================

This program generates thumbnails for the available animations found on the server by
connecting as a Websocket client and recording the incoming OPC messages.

It can then either output a single SVG file containing the recorded animation, or generate one SVG for each frame of the animation (these can be used to create animated gifs using Imagemagick).

Displaying an animated SVG as a single file is very CPU intensive in the browser, so using GIFs is better.

Requirements
------------

To be able to run this you will need the SVG [`siknas-skylt.svg`](siknas-skylt.svg) and [`layout.json`](layout.json) (which contains the positions of the LEDs).

Running
-------

```bash
docker build -t siknas-skylt-thumbgen .

# Get dependencies
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen dep ensure -v

# Help
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen go run *.go --help

# Generate a single SVG animation.
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen go run *.go --host $(docker-machine ip):8080

# Generate one frame per OPC message.
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen go run *.go --host $(docker-machine ip):8080 --output-frames
```

ImageMagick and gifs
--------------------

To create an animated gif out of the generated SVGs, [ImageMagick](https://www.imagemagick.org) can be used:

```bash
# Docker image for ImageMagick (or install it natively is preferred).
docker build -t svg2gif -f Dockerfile.svg2gif .

# Convert a set of SVG frames to an animated gif using Imagemagick (Takes a long time).
docker run -it --rm -v $(pwd):/shared -w /shared svg2gif ./makegifs.sh
```

**NOTE** On Windows running this in docker takes 2x more time!
