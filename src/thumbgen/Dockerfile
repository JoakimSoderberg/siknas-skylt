FROM node:9-alpine

RUN apk update && apk add build-base cairo-dev pango-dev giflib-dev 
RUN yarn --no-bin-links

CMD node index.js
