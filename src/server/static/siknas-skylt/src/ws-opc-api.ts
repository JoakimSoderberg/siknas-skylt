import { autoinject } from 'aurelia-framework';
import { EventAggregator } from 'aurelia-event-aggregator';
import { WebsocketOPCMessage } from "./messages";

const MAX_BACKOFF = 5000;
const BACKOFF_INCR = 500;

// TODO: Make Base class with the backoff support.
@autoinject()
export class WSOPCAPI {

    backoff: number;
    socket: WebSocket;

    constructor(private events: EventAggregator) { }

    connect() {
        try {
            this.socket = new WebSocket(`ws://${location.host}/ws/opc`);
            this.socket.onmessage = (e: MessageEvent) => { this.onmessage(e); };
            this.socket.onclose = (e: CloseEvent) => { this.onclose(e); };
            this.socket.onerror = (e: ErrorEvent) => { this.onerror(e); };
            this.socket.onopen = (e: Event) => { this.onopen(e); };
        }
        catch (ex) {
            console.log(ex);
        }
    }

    onmessage(e: MessageEvent) {
        //console.log("OPC Websocket message:", e);

        // TODO: Read binary data.
        this.events.publish(new WebsocketOPCMessage(e.data));
    }

    onerror(e: ErrorEvent) {
    }

    onopen(e: Event) {
        this.backoff = 0;
        console.log("OPC Websocket connected");
    }

    incrementBackoff() {
        this.backoff = Math.min(MAX_BACKOFF, this.backoff + BACKOFF_INCR);
    }

    onclose(e: CloseEvent) {
        console.log("OPC Websocket disconnected (reconnect in " + this.backoff + "s) " + e.reason);
        this.incrementBackoff();
        setTimeout(
            () => {
                console.log("OPC Websocket attempting reconnect...");
                this.connect();
            },
            this.backoff
        );
    }

    get isConnected(): boolean {
        if (this.socket) {
            return (this.socket.readyState == WebSocket.OPEN);
        }

        return false;
    }
}
