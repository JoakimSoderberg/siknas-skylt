#!/bin/sh
set -e

echo "=============== INSTALL YARN PACKAGES ==============="
yarn --no-bin-links

echo "=============== BUILDING AURELIA PROJECT ==============="
aurelia build

echo "=============== COPYING BUILD OUTPUT ==============="
mkdir -p build/static/siknas-skylt/
cp -R -f index.html favicon.ico scripts/ fonts/ images/ styles/ misc/ \
    build/static/siknas-skylt/
