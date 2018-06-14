#!/bin/bash

batch=20
output_path=""

show_help() {
    echo "Usage: $0 <options> [*.svg...]"
    echo ""
    echo "Make a gif for a single Animation from the generated SVG frames created by thumbgen"
    echo ""
    echo "  -h          Show this help."
    echo "  -o <path>   Change output filename for the gif (Example output/flames.gif)"
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

files=( $@ )
file_out_path="${output_path}"

base=$(basename ${file_out_path})
anim=${base%.*}

echo "Animation: ${anim} (${#files[@]} frames)"

if [ ${#files[@]} = 0 ]; then
    echo "  No frames found, skipping ..."
    return
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
