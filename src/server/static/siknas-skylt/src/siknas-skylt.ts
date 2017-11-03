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

interface OPCColor {
    r: number;
    g: number;
    b: number;
}

const OPC_HEADER_LEN = 4;

@autoinject()
export class SiknasSkylt {
    svg: any;
    skylt: any;

    layout: OPCPixel[];
    //pixels: OPCPixel[];
    pixel: any;

    constructor(private opc: WSOPCAPI, private events: EventAggregator, private client: HttpClient) {
        this.events.subscribe(WebsocketOPCMessage, msg => {
            let d = new Uint8Array(msg.data);
            let channel = d[0];
            let command = d[1];
            let highLen = d[2];
            let lowLen = d[3];
            let length = (highLen << 8) | lowLen;
            //console.log("length:" + length + " high: " + highLen + " low: " + lowLen);

            console.log("length: " + length + " length layout: " + this.layout.length);

            for (let i = OPC_HEADER_LEN, j = 0; i < length; i += 3, j++) {
                let r = d[i]
                let g = d[i + 1]
                let b = d[i + 2]

                this.layout[j].color = [r, g, b];
            }
            this.updatePixels();
        })
    }

    created() {
        this.opc.connect();
    }

    updatePixels() {
        if (!this.layout)
            return;

        let w = parseInt(this.svg.style("width"), 10);
        let h = parseInt(this.svg.style("height"), 10);

        this.pixel = this.svg.selectAll("circle").data(this.layout)
            .enter()
            .append("circle")
            .attr("cx", function (p: OPCPixel) {
                return (p.point[0] * w * 0.81) + (w * 0.05);
            })
            .attr("cy", function (p: OPCPixel) {
                return (p.point[1] * h * 0.55) + (h * 0.20);
            })
            .attr("r", 2)
            .style("fill", function (p: OPCPixel) {
                return `rgb(${p.color[0]},${p.color[1]},${p.color[2]})`
            });
        //.style("stroke", "black")
        //.style("stroke-width", 1.0);

        this.pixel.exit().remove();
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
                this.updatePixels();
            });

        this.skylt = this.svg.append("image").attr("href", "images/siknas-skylt.svg")
            .attr("x", 0)
            .attr("y", 0)
            .attr('width', 300)
            .attr('height', 300);
    }
}
