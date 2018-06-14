#!/bin/bash

set -e

# We have to batch the work or ImageMagick runs out of memory.
batch=20
output_path="output"
input_path="output"

show_help() {
    echo ""
    echo "Make gifs from the generated SVG frames created by thumbgen"
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

mkdir -p ${output_path}

# We assume all directories are animations.
animations=( $(echo ${input_path}/*/ | xargs basename) )

for anim in "${animations[@]}"
do
    ./makegif.sh -b ${batch} -o "${output_path}/${anim}.gif" ${input_path}/${anim}/*.svg
done
