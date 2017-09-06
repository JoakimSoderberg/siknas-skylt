
h1. Scripts

This directory contains various helper scripts.

To run either `pip install` on your system, or you can build the
provided docker and run in that:

```bash
# Build the docker image and install PIP requirements.
docker build -t siknas-skylt-scripts .

# Run the script.
docker run -it --rm siknas-skylt-scripts rescale_layout.py --help

# Mount a volume with any files.
docker run -it --rm -v $(pwd)/../layouts/:/layouts siknas-skylt-scripts rescale_layout.py /layouts/layout.json
```