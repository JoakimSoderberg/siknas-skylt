#!/bin/bash

# TODO: Add command line parsing --quiet flag
# TODO: Default to error if missing build

echo "#########################################################"
echo "NOTE!!! Since Processing CLI is broken you have to build"
echo "the processing animations from within the GUI."
echo
echo "For each animation do this:"
echo "  File -> Export Application"
echo "    Select    'Linux'"
echo "    Unselect  'Embed Java for Windows (64-bit)'"
echo "  Click Export"
echo
echo "#########################################################"
echo
echo "This script looks for these directories:"
echo "  animations/*/application.{amd64,armv6hf}"
echo

read -p "Ready to continue? [y/n] " answer
if [ "${answer}" = "n" ]; then
    exit 0
fi

anim_dir="./examples"
out_amd64="./build/animations-amd64"
out_armhf="./build/animations-armhf"

mkdir -p ${out_amd64}
mkdir -p ${out_armhf}

for animation in $(ls ${anim_dir})
do
    echo "Animation: ${animation}"

    curdir="${anim_dir}/${animation}"

    if [ -d "${curdir}/application.linux64/" ]; then
        echo "  [linux64] OK! Copying to build dir..."
        cp -R ${curdir}/application.linux64 ${out_amd64}/${animation}
    else
        echo "  [linux64] -"
    fi

    if [ -d "${curdir}/application.linux-armv6hf/" ]; then
        echo "  [linux-armv6hf] OK! Copying to build dir..."
        cp -R ${curdir}/application.linux-armv6hf ${out_armhf}/${animation}
    else
        echo "  [linux-armv6hf] -"
    fi  
done

# TODO: Write config that sets up to run all animations

