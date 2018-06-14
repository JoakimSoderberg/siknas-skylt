Animation sketches
==================

This folder contains a list of Animations used with the display.

To enable running **Processing** sketches in it `xvfb` is used.

To build a sketch as a standalone application:

```bash
# Assuming you're in this directory. Note the sketch path must be an absolute path.
processing-java --sketch=$(pwd)/flames/ --platform=linux --export
```

This will create a bunch of `application.*` directories for the different platform in the sketch directory.

**NOTE** Do not use `--output` if you want applications for all supported platforms!

