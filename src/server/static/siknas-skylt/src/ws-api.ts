import { autoinject } from 'aurelia-framework';
import { EventAggregator } from 'aurelia-event-aggregator';
import {
    WebsocketDisconnected, WebsocketConnected, WebsocketError,
    WebsocketMessageReceived, WebsocketAnimationList
} from "./messages";
import { Animation } from "./types";

const MAX_BACKOFF = 5000;
const BACKOFF_INCR = 500;

@autoinject()
export class WSAPI {

    backoff: number
    socket: WebSocket;
    active: boolean;

    constructor(private events: EventAggregator) { }

    connect() {
        this.socket = new WebSocket(`ws://${location.host}/ws`);
        this.socket.onmessage = (e: MessageEvent) => { this.onmessage(e); };
        this.socket.onclose = (e: CloseEvent) => { this.onclose(e); };
        this.socket.onerror = (e: ErrorEvent) => { this.onerror(e); };
        this.socket.onopen = (e: Event) => { this.onopen(e); };
    }
    
    onmessage(e: MessageEvent) {
        // Raise the raw message.
        this.events.publish(new WebsocketMessageReceived(e));

        if (e.data["message_type"] == "list") {
            // http://choly.ca/post/typescript-json/
            // TODO: Would be nice if this was typesafe instead.
            this.events.publish(new WebsocketAnimationList(e.data));
        }
    }

    onerror(e: ErrorEvent) {
        this.events.publish(new WebsocketError(e))
    }

    onopen(e: Event) {
        this.events.publish(new WebsocketConnected())
    }

    incrementBackoff() {
        this.backoff = Math.min(MAX_BACKOFF, this.backoff + BACKOFF_INCR);
    }

    onclose(e: CloseEvent) {
        this.incrementBackoff();
        setTimeout(
            () => {
                if (this.active) this.connect();
            },
            this.backoff * 1000
        );
        this.events.publish(new WebsocketDisconnected(e))
    }
}
