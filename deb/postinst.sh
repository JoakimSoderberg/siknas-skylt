#!/bin/bash

# TODO: Create custom user to run under??
# TODO: Explain what commands to run after install.

echo "Getting started"
echo "---------------"
echo
echo "0. Install Fadecandy "  # TODO: Detailed instructions...
echo
echo "1. Create a config file to get started based on the example:"
echo
echo "  sudo cp /etc/siknas/siknas-example.yaml /etc/siknas/siknas.yaml"
echo
echo "2. Verify Xvfb is installed and working:"
echo
echo "  Xvfb --help"  # TODO: Add a script that tests this.
echo
echo "3. Start the Siknas Skylt Server service:"
echo
echo "  sudo systemctl start fcserver"
echo "  sudo systemctl start Xvfb"
echo "  sudo systemctl start siknas-skylt"
echo
