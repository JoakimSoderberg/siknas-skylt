#!/bin/bash

echo "Starting Xvfb inside server container"
docker-compose exec server Xvfb :1 -screen 0, 1024x768x16
