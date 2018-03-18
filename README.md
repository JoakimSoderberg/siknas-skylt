
Siknas skylt
============

This repository contains code for running a custom interactive LED display, together with some helper tools.

The display uses the Fadecandy LED controller board https://github.com/scanlime/fadecandy that supports
the OpenPixelControl (OPC).

To use this software, a **Fadecandy** server needs to run and be connected ot the display.

Repository structure
--------------------

* [`docs`](docs/) - Contains some documentation on how the display works.
* [`examples/`](examples/) - Contains Processing sketches that animates the display using OPC.
* [`image-gui/`](image-gui/) - A .NET C# program used to map real pixel locations to the virtual ones. (Used to produce `layout.json` that the Processing sketches use).
* [`layouts`](layouts/) - Contains the [`layout.json`](layouts/layout.json) created by using the [`image-gui`](image-gui/), and the source image used to do this.
* [`scripts/`](scripts/) - A script to re-scale the coordinates in [`layout.json`](layouts/layout.json).
* [`src/controlpanel/`](src/controlpanel/) - A websocket client that talks to the control panel via a serial port over USB. The Websocket client connects to the server. Written in Golang.
* [`src/server/`](src/server/) - A server that hosts an OPC proxy, as well as a webserver and websockets server. This forwards the OPC traffic to the display coming from the processing sketches (that it starts and stops). It also broadcasts the traffic to connected webclients via websockets. The webserver hosts a web page that will let the user chose which processing sketch to run.
* [`src/server/static/`](src/server/static/) - Hosts a webpage written in [Aurelia](https://aurelia.io/) that displays the list of Processing sketches used to animate the display. This webpage also listens to the binary OPC traffic, and uses [D3](https://d3js.org/) to animate an SVG copy of the real display.

System overview
---------------

![System overview](docs/system-design/siknas-skylt.svg)