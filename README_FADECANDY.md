Fadecandy
=========

To use this project you will need [**Fadecandy**](https://github.com/scanlime/fadecandy) installed somewhere.

Most likely on the same RPi, but this is not a requirement as long as **Siknas-skylt server** has a network connection to it.

The official releases do not contain a debian package:
https://github.com/scanlime/fadecandy/releases

But it works fine to use those. However we want to install it via the debian package which also includes a **SystemD service unit**.

# Alternative 1
To build a debian package yourself on the Raspberry pi:

```bash
sudo apt-get install git cmake build-essential

git clone git@github.com:scanlime/fadecandy.git
cd fadecandy

git submodule update --init

cd server
mkdir build
cd build

cmake ..
make
make package  # For debian package.

# Install it (and fix any broken dependencies).
sudo dpkg -i fcserver*.deb
sudo apt-get -f install

sudo systemctl start fcserver
```

# Alternative 2

Install the official release package (https://github.com/scanlime/fadecandy/releases) and create the **SystemD service yourself**:

`/lib/systemd/system/fcserver.service`
```ini
[Unit]
Description=Fadecandy USB LED controller server

[Service]
ExecStart=/usr/local/bin/fcserver
RemainAfterExit=yes
StandardOutput=journal+console
StandardError=journal+console

[Install]
WantedBy=multi-user.target
```


