import { autoinject } from 'aurelia-framework';
import * as d3 from "d3";
import { WSOPCAPI } from "./ws-opc-api";
import { EventAggregator } from 'aurelia-event-aggregator';
import { WebsocketOPCMessage } from "./messages";
import { HttpClient } from 'aurelia-fetch-client';

type OPCPoint = [number, number, number]

interface OPCPixel {
    point: OPCPoint;
    color: OPCPoint;
}

const OPC_HEADER_LEN = 4;

@autoinject()
export class SiknasSkylt {
    svg: any;
    skylt: any;
    // TODO: make it possible to disable the animation via a bindable property.

    layout: OPCPixel[];
    pixelsCreated: boolean = false;
    pixels: any;

    constructor(private opc: WSOPCAPI, private events: EventAggregator, private client: HttpClient) {
        this.events.subscribe(WebsocketOPCMessage, msg => {
            if (!this.pixelsCreated)
                return;

            let d = new Uint8Array(msg.data);
            let channel = d[0];
            let command = d[1];
            let highLen = d[2];
            let lowLen = d[3];
            let length = (highLen << 8) | lowLen;

            // TODO: Would be good if we could set the style directly instead of first copying here...
            for (let i = OPC_HEADER_LEN, j = 0; i < length; i += 3, j++) {
                let r = d[i]
                let g = d[i + 1]
                let b = d[i + 2]

                this.layout[j].color = [r, g, b];
            }

            // Change color for all pixels.
            this.pixels.style("fill", function (p: OPCPixel) {
                return `rgb(${p.color[0]},${p.color[1]},${p.color[2]})`
            });
        })
    }

    created() {
        this.opc.connect();
    }

    createPixels() {
        if (!this.layout)
            return;

        let w = parseInt(this.svg.style("width"), 10);
        let h = parseInt(this.svg.style("height"), 10);

        var defs = this.svg.append("defs");
        var filter = defs.append("filter")
            .attr("id", "glow")
            .append("feGaussianBlur")
            .attr("stdDeviation", "0.5")
            .attr("result", "coloredBlur");

        // Merge the original shape with the blur.
        var feMerge = filter.append("feMerge");
        feMerge.append("feMergeNode")
            .attr("in", "coloredBlur");
        feMerge.append("feMergeNode")
            .attr("in", "SourceGraphic");

        this.pixels = this.svg.selectAll("circle").data(this.layout)
            .enter()
            .append("circle")
            // These values are hand tweaked to be placed over the logo properly.
            .attr("cx", function (p: OPCPixel) {
                return (p.point[0] * w * 0.81) + (w * 0.05);
            })
            .attr("cy", function (p: OPCPixel) {
                return (p.point[1] * h * 0.55) + (h * 0.20);
            })
            .attr("r", 2)
            .attr("filter", "url(#glow)")
            .style("fill", function (p: OPCPixel) {
                return `rgb(${p.color[0]},${p.color[1]},${p.color[2]})`
            });

        this.pixels.exit().remove();

        this.pixelsCreated = true;
    }

    attached() {
        this.svg = d3.select("#siknas-skylt").append('svg')
            .attr('width', 300)
            .attr('height', 300);

        this.client.fetch("misc/layout.json")
            .then(response => response.json())
            .then(data => {
                this.layout = data;

                for (let p of this.layout) {
                    p.color = [0, 0, 0];
                }
                console.log("Got layout: ", this.layout);
                this.createPixels();
            });

        this.skylt = this.svg.append("image").attr("href", "images/siknas-skylt.svg")
            .attr("x", 0)
            .attr("y", 0)
            .attr('width', 300)
            .attr('height', 300);
    }
}
