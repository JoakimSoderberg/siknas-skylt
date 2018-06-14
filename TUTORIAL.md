Tutorial (Development)
======================

This tutorial uses **Docker**.

1. Install Docker. (On Windows use [**Docker toolbox**](https://docs.docker.com/toolbox/toolbox_install_windows/)).

2. Start the server. (On Windows run in the **Docker console**):

    ```bash
    docker-compose up -d
    docker logs -f server   # Depending what you want to follow.
    docker logs -f aurelia

    # Get error on starting aurelia?
    cd src/server/static/siknas-skylt/node_modules/
    rm -rf npm
    docker-compose up -d aurelia
    ```

3. Now you can surf to the server: http://localhost:8080 (Use `docker-machine ip` to get IP on Windows instead of `localhost`)
    ```bash
    start http://$(docker-machine ip):8080  # To open the webpage from the console on Windows.
    ```
    Which looks like this:

    ![Siknäs skylt website](docs/images/website.png)

4. Download and install [**Processing**](https://processing.org/).

    Make sure the installation directory is added to the **Path** in Windows (Or write the full path in the examples below for `processing-java`).

5. Open the [**flames**](examples/flames/flames.pde) example sketch.

    **NOTE** If using the **Processing GUI** on Windows you must edit the `host` variable at the top of the sketch to connect to the `docker-machine ip` instead of `localhost`

    ```bash
    # Assuming you are standing in the root of this repository.
    processing-java --help

    # To run the example sketch (Note, must be the full path)
    processing-java --sketch=$(pwd)/examples/flames/ --run

    # On Windows you have to specify the OPC host to connect to since it is not localhost.
    processing-java --sketch=$(pwd)/examples/flames/ --run $(docker-machine ip):7890

    # NOTE if you keep getting disconnected repeatedly. Make sure you connect to port 7890.
    ```

6. Now if you don't have a real display built yet, you can test with the [**Simulator**](https://github.com/JoakimSoderberg/OPCSim). Download the latest version here: https://github.com/JoakimSoderberg/OPCSim/releases/latest

    Once you have the simulator running, you must configure the server so that it knows how to connect to it. To do this there's a [config file example](src/server/siknas.yaml.example).

    The server will look for this file in either the same path, or under `/etc/siknas/` on the system.

    So create a copy of it:

    ```bash
    cd src/server/
    cp siknas.yaml.example siknas.yaml
    cat siknas.yaml
    cd ../..
    ```

    Example contents:
    ```yaml
    opc-servers:
        simulator:
            host: 192.168.1.75
            port: 7890
    #processes:
    #    Flames:
    #        description: Cool flames
    #        Exec: /path/to/sketch/executable
    #        KillCommand: kill -9 `ps aux | grep path/sketch` | awk '{print $1}'
    ```

    Edit this to connect to the correct ip. (On **Windows** this is your docker host IP, see `ipconfig`, on **Linux** most likely `localhost`).

    You can make the server forward the traffic both towards the real display, as well as the simualtor by adding a second entry in the config.

    Now restart the server to reload the config:

    ```bash
    docker-compose restart server
    ```

7. With the above steps, running **Processing** manually means we are not using the webserver to choose the animation. To enable that we can add a list of processes in the config file:

    Since this will run inside of docker you'll need to export the Processing sketch you want to run.

    ```bash
    processing-java --sketch=$(pwd)/examples/flames/ --platform=linux --export
    ```

    The [`docker-compose.yml`](docker-compose.yml) project will mount the `/examples` directory under `/animations`

    ```yaml
    # ...
    processes:
        Flames:
            description: Cool flames
            Exec: /animations/flames/application.linux64/flames
            # The kill command is required since a normal kill won't kill the Java process otherwise.
            KillCommand: kill -9 `ps aux | grep animations/flames | awk '{print $1}'`
    ```

    Restart the server so the config is reloaded:

    ```bash
    docker-compose restart server
    ```

    It won't work yet. The processing sketch must run in **headless mode**. To do that we must use `Xvfb` (a virtual X server), in production we can set this up at system startup. You will see errors like this:

    ```bash
    Animation process died unexpectedly restarting: exit status 1
    ```

    When running in docker we need to start it manually (before attempting to start any process).

    ```bash
    docker-compose exec server Xvfb :1 -screen 0, 1024x768x16 &
    ```

8. Surf to http://localhost:8080 and select the animation and press play. The logo should start animating.

