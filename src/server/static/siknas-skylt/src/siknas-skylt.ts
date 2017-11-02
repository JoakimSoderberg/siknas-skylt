import { autoinject } from 'aurelia-framework';
import * as d3 from "d3";
import { WSOPCAPI } from "./ws-opc-api";
import { EventAggregator } from 'aurelia-event-aggregator';
import { WebsocketOPCMessage } from "./messages";

@autoinject()
export class SiknasSkylt {
    svg: any;

    constructor(private opc: WSOPCAPI, private events: EventAggregator) {
        this.events.subscribe(WebsocketOPCMessage, msg => {
            //console.log("OPC msg:", msg);
        })
    }

    created() {
        this.opc.connect();
    }

    attached() {
        this.svg = d3.select("#siknas-skylt").append('svg')
            .attr('width', 300)
            .attr('height', 300);

        this.svg.append("image").attr("href", "images/siknas-skylt.svg")
            .attr("x", 0)
            .attr("y", 0)
            .attr('width', 150)
            .attr('height', 150);
    }
}
