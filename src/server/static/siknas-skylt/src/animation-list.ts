import { EventAggregator } from 'aurelia-event-aggregator';
import {
    WebsocketConnected, WebsocketDisconnected, WebsocketMessageReceived,
    WebsocketError, WebsocketAnimationList, AnimationViewed
} from './messages';
import { autoinject } from 'aurelia-framework';
import { Animation, AnimationListMessage } from "./types";
import { WSAPI } from "./ws-api";

@autoinject()
export class AnimationList {
    animations: Array<Animation>;
    playingAnimation: Animation | null;

    constructor(private events: EventAggregator, private api: WSAPI) {
        // Once the websocket is connected we'll receive a list of animations.
        events.subscribe(WebsocketAnimationList, msg_raw => {
            let msg: AnimationListMessage = msg_raw.data;
            console.log("Animations received:", msg);

            this.animations = msg.anims;

            if (msg.playing > 0) {
                this.playingAnimation = msg.anims[msg.playing];
            } else {
                this.playingAnimation = null;
            }
        });
    }

    play(animation: Animation | null) {
        if (animation != null) {
            this.api.sendPlayMessage(animation.name);
        }

        return true;
    }

    stop() {
        this.api.sendPlayMessage("");
    }
}
