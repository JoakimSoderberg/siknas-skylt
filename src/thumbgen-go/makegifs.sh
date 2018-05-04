#!/bin/bash

set -e

# We have to batch the work or ImageMagick runs out of memory.
batch=20

# We assume all directories are animations.
animations=$(ls -d output/*/ | xargs basename)

for anim in "${animations[@]}"
do
    files=( output/flames/*.svg )
    echo "Animation: ${anim} (${#files[@]} frames)"

    if [ ${#files[@]} = 0 ]; then
        echo "  No frames found, skipping ..."
        continue
    fi

    mkdir -p tmp/
    rm -rf tmp/${anim}.*.miff

    # Generate intermediate batches of animations in the internal ImageMagick miff format.
    for (( i=0; $i<${#files[@]}; i+=$batch ))
    do
        intname="tmp/${anim}.$(printf "%06d" $i).miff"
        echo "  Generating ($i of ${#files[@]}) ${intname}"
        magick -delay 2 -loop 0 "${files[@]:$i:$batch}" -scale 150x150 ${intname}
    done

    echo -n
    echo "  Combining result for ${anim} into output/${anim}.gif"
    magick tmp/${anim}.*.miff output/${anim}.gif
done
