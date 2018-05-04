#!/bin/bash

# TODO: Remove hard coded names
files=( output/flames/*.svg )
batch=20

mkdir -p output/tmp/

for (( i=0; $i<${#files[@]}; i+=$batch ))
do
    echo "Generating output/tmp/flames.$(printf "%06d" $i).gif ..."
    convert -delay 2 -loop 0 "${files[@]:$i:$batch}" -scale 150x150 output/tmp/flames.$(printf "%06d" $i).gif
done

echo "Combining result..."
convert output/tmp/flames.*.gif output/flames.gif
