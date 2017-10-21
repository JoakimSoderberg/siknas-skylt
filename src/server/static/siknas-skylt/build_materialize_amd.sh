#!/bin/sh

cp ./node_modules/aurelia-materialize-bridge/build/tools/*.js .
./node_modules/.bin/r.js -o rbuild.js
