Siknäs-skylt Thumbnail Generator
================================

This program generates thumbnails for the available animations found on the server by
connecting as a Websocket client and recording the incoming OPC messages.

It can then either output a single SVG file containing the recorded animation, or generate one SVG for each frame of the animation (these can be used to create animated gifs using Imagemagick).

Displaying an animated SVG as a single file is very CPU intensive in the browser, so using GIFs is better.

Running
-------

```bash
docker build -t siknas-skylt-thumbgen-go .

# Get dependencies
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen-go dep ensure -v

# Help
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen-go go run *.go --help

# Generate a single SVG animation.
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen-go go run *.go --host $(docker-machine ip):8080

# Generate one frame per OPC message.
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen-go go run *.go --host $(docker-machine ip):8080 --output-frames 
```

Docker (rsvg / Imagemagick)
---------------------------

To rasterize a single SVG `rsvg-convert` can be used:

```bash
docker build -t svg2gif -f Dockerfile.svg2gif .

# Convert a single SVG to png (animations not supported) using librsvg.
docker run -it --rm -v $(pwd):/shared -w /shared svg2gif rsvg-convert some.svg -o some.png

# Convert a set of SVG frames to an animated gif using Imagemagick.
docker run -it --rm -v $(pwd):/shared -w /shared svg2gif ./gifmaker.sh
```


