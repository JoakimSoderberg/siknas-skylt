#!/usr/bin/env python

from __future__ import print_function
import argparse
import json
import os
import sys
from PIL import Image

verbose = False


def eprint(*args, **kwargs):
    print(*args, file=sys.stderr, **kwargs)


def dprint(*args, **kwargs):
    if verbose:
        print(*args, **kwargs)


def main():
    parser = argparse.ArgumentParser(description='Rescales a 2D Openpixel protocol layout file, based on the original '
                                     ' image size, assuming the coordinates are pixel coordinates of that image file.')

    parser.add_argument('layout', help='Path to layout JSON file.')
    parser.add_argument('--image', help='Path to the image the layout was created from.')
    parser.add_argument('--size', nargs=2, type=int, help='If no image is available give the size manually.')
    parser.add_argument('--scale', '-s', default=1.0, type=float, help='Scale points by this factor after normalizing.')
    parser.add_argument('--output', '-o', help='Output layout JSON file. Prints to stdout if not specified.')
    parser.add_argument('--verbose', '-v', action='store_true', help='Verbose output.')

    args = parser.parse_args()

    dprint("Opening {}".format(args.layout))

    new_points = []
    img_size = (1, 1)
    global verbose
    verbose = args.verbose

    with open(args.layout) as f:
        points = json.load(f)

        if args.image:
            with Image.open(args.image) as imf:
                img_size = imf.size
        elif args.size:
            img_size = tuple(args.size)

        for p in points:
            p = p['point']
            org_x = p[0]
            org_y = p[1]
            x = org_x / float(img_size[0]) * args.scale
            y = org_y / float(img_size[1]) * args.scale
            dprint('{}, {}  => {}, {}'.format(org_x, org_y, x, y))
            new_points.append({'point': [x, y, 0]})

    if args.output:
        with open(args.output, 'w') as f:
            json.dump(new_points, f, indent=4)
            dprint('Wrote output to {}'.format(args.output))
    else:
        print(json.dumps(new_points, indent=4))

    dprint('Rescaled points using image size {} and scale factor {}'.format(img_size, args.scale))


if __name__ == '__main__':
    try:
        main()
    except Exception as ex:
        eprint('Error: {}'.format(ex))
