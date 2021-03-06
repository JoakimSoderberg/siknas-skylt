import { autoinject, bindable } from 'aurelia-framework';
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
    @bindable disableAnimation: boolean = false;

    layout: OPCPixel[];
    pixelsCreated: boolean = false;
    pixels: any;

    @bindable blurlevel: number;

    constructor(private opc: WSOPCAPI, private events: EventAggregator, private client: HttpClient) {
        this.events.subscribe(WebsocketOPCMessage, msg => {
            if (!this.pixelsCreated || this.disableAnimation)
                return;

            let d = new Uint8Array(msg.data);
            let channel = d[0];
            let command = d[1];
            let highLen = d[2];
            let lowLen = d[3];
            let length = (highLen << 8) | lowLen;

            if (command == 0) {
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
            } else if (command == 0xff) {
                // System exclusive command.
                // We are only looking for Fadecandy ColorCorrection messages.
                let highSystemID = d[OPC_HEADER_LEN];
                let lowSystemID = d[OPC_HEADER_LEN + 1];
                let systemID = (highSystemID << 8) | lowSystemID;

                if (systemID != 1) {
                    // Not Fadecandy
                    return;
                }

                let highCommandID = d[OPC_HEADER_LEN + 2];
                let lowCommandID = d[OPC_HEADER_LEN + 3];
                let commandID = (highCommandID << 8) | lowCommandID;

                if (commandID != 1) {
                    // Not a color correction Fadecandy command.
                    return;
                }

                let brightness = 0;

                // We subtract 4 from length to exclude the systemID + commandID.
                let jsonAsBytes = d.slice(OPC_HEADER_LEN + 4);
                var jsonTxt = String.fromCharCode.apply(null, jsonAsBytes);
                let j = JSON.parse(jsonTxt);

                // TODO: Remove debug print
                // console.log("Brightness: ", j["whitepoint"][0]);

                let red = j["whitepoint"][0];
                let green = j["whitepoint"][1];
                let blue = j["whitepoint"][2];

                document.getElementById("colorCorrectMatrix").setAttribute("values",
                    `${red} 0        0       0    0
                     0      ${green} 0       0    0
                     0      0        ${blue} 0    0
                     0      0        0       1    0 `)
            }
        })
    }

    created() {
        this.opc.connect();
    }

    createPixels() {
        if (!this.layout)
            return;

        let viewBox = this.svg.attr("viewBox").split(" ");
        let size = viewBox.slice(2);
        let w = size[0];
        let h = size[1];
        let aspect = w / h;
        console.log("w: " + w + " h: " + h);

        let defs = this.svg.append("defs");
        let filter = defs.append("filter");

        // We use this for the color correction packets
        // (Which is used to set brightness).
        filter.attr("id", "siknasFilter")
            .append("feColorMatrix")
            .attr("id", "colorCorrectMatrix")
            .attr("type", "matrix")
            .attr("values",
                `1   0   0   0   0
                 0   1   0   0   0
                 0   0   1   0   0
                 0   0   0   1   0`);

        // Blur to look like the real display.
        filter.append("feGaussianBlur")
            .attr("id", "blurFilter")
            .attr("stdDeviation", this.blurlevel.toString());

        // Change the color of the Siknäs text of the SVG.
        this.svg.select("#Siknas").selectAll("path").style("fill", "black");

        let grp = this.svg.append("g")
            // A copy of the "Siknäs" text paths exists in the SVG defined as a clip-path
            // we want the pixels to be stay inside of this.
            .attr("clip-path", "url(#SiknasClipPath)")
            .attr("style", "filter: url(#siknasFilter)");

        this.pixels = grp.selectAll("circle").data(this.layout)
            .enter()
            .append("circle")
            // These values are hand tweaked to be placed over the logo properly.
            .attr("cx", function (p: OPCPixel) {
                return (p.point[0] * w * 0.81) + (w * 0.05);
            })
            .attr("cy", function (p: OPCPixel) {
                return (p.point[1] * h * 0.55) + (h * 0.20);
            })
            .attr("r", 10)
            .style("fill", function (p: OPCPixel) {
                return `rgb(${p.color[0]},${p.color[1]},${p.color[2]})`
            });

        this.pixels.exit().remove();
        this.pixelsCreated = true;
    }

    attached() {
        // Read the SVG as XML instead of adding it as an <image>
        // so we can manipulate parts of it directly.
        this.client.fetch("images/siknas-skylt.svg")
            .then(response => response.text())
            .then(str => (new DOMParser()).parseFromString(str, "image/svg+xml"))
            .then(data => {
                // Add SVG to DOM.
                document.getElementById("siknas-skylt")
                    .appendChild(data.documentElement);

                this.svg = d3.select("#siknas-skylt svg")
                    .attr('width', 300)
                    .attr('height', 300);

                // Get the Pixel layout and add them to the SVG.
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
            });
    }

    blurlevelChanged(newValue: number, oldValue: number) {
        if (this.svg) {
            this.svg.select("#blurFilter").attr("stdDeviation", newValue.toString());
        }
    }
}
