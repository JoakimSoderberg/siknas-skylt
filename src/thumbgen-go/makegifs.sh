#!/bin/bash

set -e

# We have to batch the work or ImageMagick runs out of memory.
batch=20

# We assume all directories are animations.
animations=$(ls -d output/*/ | xargs basename)

for anim in "${animations[@]}"
do
    files=( output/${anim}/*.svg )
    out_path="output/${anim}.gif"

    echo "Animation: ${anim} (${#files[@]} frames)"

    if [ ${#files[@]} = 0 ]; then
        echo "  No frames found, skipping ..."
        continue
    fi

    mkdir -p tmp/
    rm -rf tmp/${anim}.*.miff

    if [ ${batch} -ge ${#files[@]} ]; then
        echo "  Generating ${anim} gif ${out_path}"
        magick -delay 2 -loop 0 -background none "${files[@]}" -scale 150x150 ${out_path}
    else
        # Generate intermediate batches of animations in the internal ImageMagick miff format.
        for (( i=0; $i<${#files[@]}; i+=$batch ))
        do
            intermediate_path="tmp/${anim}.$(printf "%06d" $i).miff"

            echo "  Generating ($i of ${#files[@]}) ${intermediate_path}"
            magick -delay 2 -loop 0 -background none "${files[@]:$i:$batch}" -scale 150x150 ${intermediate_path}
        done

        echo -n
        echo "  Combining result for ${anim} into ${out_path}"
        magick -background none tmp/${anim}.*.miff ${out_path}
    fi
done
