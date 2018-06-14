
Siknas skylt
============

This repository contains code for running a custom interactive LED display, together with some helper tools.

The display uses the Fadecandy LED controller board https://github.com/scanlime/fadecandy that supports
the OpenPixelControl (OPC).

To use this software, a **Fadecandy** server needs to run and be connected to the display via **Fadecandy USB controller boards**.

Build release package
--------------------

This program was designed to run on a Raspberry Pi along with **Fadcandy** fcserver.

## Creating the debian package

Currently this is a multi-stage process.

1. Build the animations themselves (Note requires using the Processing GUI, since the CLI is broken, see https://github.com/processing/processing/issues/5468).
    - Open each Processing sketch under `animations/`
    - File -> Export Application...
    - Click Export
    - Make sure the sketch dir contains `animations/animation_name/`:
        - `application.linux-armv6hf`
        - `application.linux-amd64`
2. Now run:
    ```bash
    ./build-animations.sh
    ```
3. Create thumbnails. (**Note** this will start the server and send actual traffic that is recorded and made into animated Gifs. So don't touch the webapp during this. Details [`src/thumbgen/README.md`](src/thumbgen)):
    ```bash
    ./make-thumbnails.sh
    ```
4. Now build the executables and static files (This includes Windows, OSX, Linux, Linux ARM):
    ```bash
    ./build.sh
    ```
5. Finally build the debian packages:
    ```bash
    ./build-debian.sh
    ```

## Install debian package on the Raspberry Pi

The assumption the RPi will run in "headless" mode. This means we need to use **Xvfb** for a virtual screen buffer since Processing needs that to generate the animations.

### Pre-requisite (Fadecandy)

Follow the instructions in [README_FADECANDY.md](README_FADECANDY.md) on how to build and install **Fadecandy**.

### Install

1. Copy the debian package to the RPi somehow. For example using SCP:
    ```bash
    scp build/siknas-skylt-server-armhf.deb my-rpi:.   # my-rpi is the ip or hostname of the RPi
    ```
2. Install the debian package:
    ```bash
    sudo dpkg -i ./siknas-skylt-server.armhf.deb  # This complains about unment dependencies
    sudo apt-get -f install  # Fixes dependencies.
    ```

Quickstart (Development)
------------------------

Running while developing under docker.

See the [`TUTORIAL.md`](TUTORIAL.md) for details.

```bash
# Create and edit src/server/sikas.yaml
cp src/server/siknas.yaml.example src/server/siknas.yaml

# Start server.
docker-compose up -d

# (Separate window) Run Xvfb inside of server docker.
./run_xvfb.sh

# Surf to http://localhost:8080 (Linux)
open http://$(docker-machine ip):8080   # OSX
start http://$(docker-machine ip):8080  # Windows.
```

Display
-------

Example image of the real world display.

![Siknäs skylt](docs/images/siknas-skylt.jpg)

Simulator
---------

To enable development and testing animations on the display a simulator was created in [Unity](https://unity3d.com/):
https://github.com/JoakimSoderberg/OPCSim

![Siknäs skylt simulator](docs/images/simulator.png)

Repository structure
--------------------

* [`docs`](docs/) - Contains some documentation on how the display works.
* [`animations/`](animations/) - Contains Processing sketches that animates the display using OPC.
* [`image-gui/`](image-gui/) - A .NET C# program used to map real pixel locations to the virtual ones. (Used to produce `layout.json` that the Processing sketches use).
* [`layouts`](layouts/) - Contains the [`layout.json`](layouts/layout.json) created by using the [`image-gui`](image-gui/), and the source image used to do this.
* [`scripts/`](scripts/) - A script to re-scale the coordinates in [`layout.json`](layouts/layout.json).
* [`src/controlpanel/`](src/controlpanel/) - A websocket client that talks to the control panel via a serial port over USB. The Websocket client connects to the server. Written in Golang.
* [`src/server/`](src/server/) - A server that hosts an OPC proxy, as well as a webserver and websockets server. This forwards the OPC traffic to the display coming from the processing sketches (that it starts and stops). It also broadcasts the traffic to connected webclients via websockets. The webserver hosts a web page that will let the user chose which processing sketch to run.
* [`src/server/static/`](src/server/static/) - Hosts a webpage written in [Aurelia](https://aurelia.io/) that displays the list of Processing sketches used to animate the display. This webpage also listens to the binary OPC traffic, and uses [D3](https://d3js.org/) to animate an SVG copy of the real display.
* [`src/thumbgen/`](src/thumbgen/) - Used to generate gif thumbnails by looping through and recording all the animations. These are used in the webapp to show a preview of the animation.

Generating thumbnails
---------------------

The webapp displays animated gifs as thumbnails for the animations.

To regenerate these, run the script below:

```bash
./make_thumbnails.sh
```

Or for details look at the [`README.md`](src/thumbgen) for **Thumbgen**.

System overview
---------------

![System overview](docs/system-design/siknas-skylt.svg)
