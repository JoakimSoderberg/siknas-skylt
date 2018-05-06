#!/bin/bash

set -e

# We have to batch the work or ImageMagick runs out of memory.
batch=20
output_path="output"
input_path="output"

show_help() {
    echo ""
    echo "Make gifs from the generated SVG frames"
    echo ""
    echo "  -h          Show this help."
    echo "  -i <path>   Input path, each sub-directory is expected to contain animation SVG frames."
    echo "              (default: '${input_path}')."
    echo "  -o <path>   Change output path (default: '${output_path}')."
    echo "  -b <count>  Set the batch size (default: '${batch}' frames)."
    echo "              (If this is too high ImageMagick might run out of memory)."
    echo -n
}

OPTIND=1 # Reset in case getopts has been used previously in the shell.
while getopts "h?o:i:b:" opt; do
    case "$opt" in
    h|\?)
        show_help
        exit 0
        ;;
    o)  output_path=$OPTARG
        ;;
    i)  input_path=$OPTARG
        ;;
    b)  batch=$OPTARG
        ;;
    esac
done

shift $((OPTIND-1))

echo "Batch size: ${batch} frames"
echo ""

mkdir -p ${output_path}

# We assume all directories are animations.
animations=$(ls -d ${input_path}/*/ | xargs basename)

for anim in "${animations[@]}"
do
    files=( ${input_path}/${anim}/*.svg )
    file_out_path="${output_path}/${anim}.gif"

    echo "Animation: ${anim} (${#files[@]} frames)"

    if [ ${#files[@]} = 0 ]; then
        echo "  No frames found, skipping ..."
        continue
    fi

    mkdir -p tmp/
    rm -rf tmp/${anim}.*.miff

    if [ ${batch} -ge ${#files[@]} ]; then
        echo "  Generating ${anim} gif ${file_out_path}"
        magick -delay 2 -loop 0 -background none "${files[@]}" -scale 150x150 ${file_out_path}
    else
        # Generate intermediate batches of animations in the internal ImageMagick miff format.
        for (( i=0; $i<${#files[@]}; i+=$batch ))
        do
            intermediate_path="tmp/${anim}.$(printf "%06d" $i).miff"

            echo "  Generating ($i of ${#files[@]}) ${intermediate_path}"
            magick -delay 2 -loop 0 -background none "${files[@]:$i:$batch}" -scale 150x150 ${intermediate_path}
        done

        echo ""
        echo "  Combining result for ${anim} into ${file_out_path}"
        magick -background none tmp/${anim}.*.miff ${file_out_path}
    fi
done
