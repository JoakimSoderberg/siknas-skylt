FROM node:9

RUN apt-get update \
    && apt-get install -qq libcairo2-dev libjpeg-dev libpango1.0-dev libgif-dev build-essential g++

RUN mkdir -p /opt/node/js \
    && cd /opt/node \
    npm install

COPY . /opt/node/js
WORKDIR /opt/node/js

ENTRYPOINT ["node"]
