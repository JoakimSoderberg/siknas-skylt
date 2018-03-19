Siknäs-skylt Server
===================

This server runs multiple parts:

- An OpenPixelControl (OPC) protocol proxy server. Forwards OPC traffic via both TCP and Websockets.
- A webserver hosting an Aurelia based webpage that connects via 3 different Websockets to the server.
  One for the OPC data, which it shows animated in an SVG version of the sign. One for showing the control panel
  status, and one for choosing the animation forwarded to the display.
- A Websocket that presents a list of Animations clients can choose from.
- Listens to incoming Websocket connections from the control panel, and forwards the state of it to listening websocket clients.

Docker for development
----------------------

To run on x86 using Docker:

```bash
docker build -t siknas-skylt-server .

# Ensure we have the needed dependencies on the docker host.
# (We need this when using -v during development to automatically recompile on file changes).
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-server dep ensure

# gin serves a proxy on port 3000, and recompiles on file changes.
docker run -it --rm -p 3000:3000 -v $(pwd):/go/src/app siknas-skylt-server
```

Docker on Rasberry Pi
---------------------

This setup is not intended to be used during development, but rather to run in production
on a Raspberry Pi.

```bash
docker build -t siknas-skylt-server-rpi -f Dockerfile.rpi .
docker run -it --rm -p 80:8080 siknas-skylt-server-rpi
```

Compiling exectuable
--------------------

Instead of compiling automatically inside of the docker during development:

```bash
docker build -t siknas-skylt-server .
docker run -it --rm siknas-skylt-server \
    go build -o siknas-skylt-server -v

# If you're on Windows you need to cross compile.
docker run -it --rm -e GOOS=windows siknas-skylt-server \
    go build -o siknas-skylt-server -v

# On OSX.
docker run -it --rm -e GOOS=darwin siknas-skylt-server \
    go build -o siknas-skylt-server -v

# Cross compile for Raspberry Pi that has a ARMv5 CPU.
docker run -it --rm -e GOOS=linux -e GOARCH=arm -e GOARM=5 siknas-skylt-server \
    go build -o siknas-skylt-server -v
```
