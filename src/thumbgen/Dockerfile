#FROM node:9-alpine

#RUN apk update && apk add build-base cairo-dev pango-dev giflib-dev 
#RUN yarn --no-bin-links

#CMD node index.js

FROM node:9

RUN apt-get update \
    && apt-get install -qq libcairo2-dev libjpeg-dev libpango1.0-dev libgif-dev build-essential g++

RUN mkdir -p /opt/node/js \
    && cd /opt/node \
    npm install

COPY . /opt/node/js
#    && npm i node-gyp \
#    && npm i canvas@1.6.10 d3-node@1.1.3 d3@5.0.0

WORKDIR /opt/node/js

ENTRYPOINT ["node"]