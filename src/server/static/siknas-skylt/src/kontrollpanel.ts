import { autoinject } from "aurelia-framework";
import { EventAggregator } from 'aurelia-event-aggregator';
import { HttpClient } from 'aurelia-fetch-client';
import { WebsocketControlPanelMessage } from "./messages";
import { ControlPanelMessage } from "./types";
import * as d3 from "d3";

@autoinject
export class Kontrollpanel {

    state: ControlPanelMessage;
    svg: any;

    constructor(private events: EventAggregator, private client: HttpClient) {
        this.events.subscribe(WebsocketControlPanelMessage, msg => {
            this.state = msg.data;
            console.log("kontrollpanel has msg: ", this.state, msg);
            // TODO: Animate SVG.
        });
    }

    attached() {
        this.client.fetch("images/kontrollpanel.svg")
            .then(response => response.text())
            .then(str => (new DOMParser()).parseFromString(str, "image/svg+xml"))
            .then(data => {
                // Add SVG to DOM.
                document.getElementById("kontrollpanel")
                    .appendChild(data.documentElement);
                
                this.svg = d3.select("#kontrollpanel svg")
                    .attr('width', 200);
            });
    }
}
