Control panel client
====================

This is the code for the control panel part of the display.

This will listen via the serial port (connected to the Teensy in the control panel via USB), and forward the commands sent from it over a websocket connection to the server.

The reason the server isn't directly talking via the serial port is to enable simulating the system easier. And also makes it possible for the control panel to be connected to a separate computer.

Firmware for Teensy
-------------------

The [`firmware`](firmware/) directory contains the Teensy 2.0 code that reads the analog inputs 
and sends their state over the serial port.

Quickstart
----------

```bash
docker build -t siknas-skylt-controlpanel .
docker run -it --rm siknas-skylt-controlpanel

# If you want to develop.
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-controlpanel dep ensure -v
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-controlpanel
```
