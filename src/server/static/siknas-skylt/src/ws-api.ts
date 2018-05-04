import { autoinject } from 'aurelia-framework';
import { EventAggregator } from 'aurelia-event-aggregator';
import {
    WebsocketDisconnected, WebsocketConnected, WebsocketError,
    WebsocketMessageReceived, WebsocketAnimationList, WebsocketControlPanelMessage
} from "./messages";
import { Animation, AnimationListMessage, ControlPanelMessage } from "./types";

const MAX_BACKOFF = 5000;
const BACKOFF_INCR = 500;

@autoinject()
export class WSAPI {

    backoff: number
    socket: WebSocket;

    constructor(private events: EventAggregator) { }

    connect() {
        try {
            this.socket = new WebSocket(`ws://${location.host}/ws`);
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
        console.log("Websocket message:", e);
        this.events.publish(new WebsocketMessageReceived(e));

        // http://choly.ca/post/typescript-json/
        // TODO: Would be nice if this was typesafe instead.
        let data = JSON.parse(e.data);

        switch (data["message_type"]) {
            case "list":
                let animations: AnimationListMessage = data;
                this.events.publish(new WebsocketAnimationList(animations));
                break;
            case "status":
                console.log("Status: ", data["text"]);
                break
            case "control_panel":
                let msg: ControlPanelMessage = data;
                this.events.publish(new WebsocketControlPanelMessage(data));
                break;
            default:
                console.log("Unknown message:", e.data);
                break;
        }
    }

    onerror(e: ErrorEvent) {
        this.events.publish(new WebsocketError(e))
    }

    onopen(e: Event) {
        this.events.publish(new WebsocketConnected())
        this.backoff = 0;
        console.log("Websocket connected");
    }

    incrementBackoff() {
        this.backoff = Math.min(MAX_BACKOFF, this.backoff + BACKOFF_INCR);
    }

    onclose(e: CloseEvent) {
        console.log("Websocket disconnected (reconnect in " + this.backoff + "s) " + e.reason);
        this.incrementBackoff();
        setTimeout(
            () => {
                console.log("Websocket attempting reconnect...");
                this.connect();
            },
            this.backoff
        );
        this.events.publish(new WebsocketDisconnected(e))
    }

    sendJSONMessage(msg: Object) {
        if (!this.socket)
            return;

        this.socket.send(JSON.stringify(msg));
    }

    sendSelectMessage(name: string) {
        this.sendJSONMessage({
            "message_type": "select",
            "selected": name
        });
    }

    get isConnected(): boolean {
        if (this.socket) {
            return (this.socket.readyState == WebSocket.OPEN);
        }

        return false;
    }
}
