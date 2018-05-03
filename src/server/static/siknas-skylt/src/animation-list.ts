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
        events.subscribe(WebsocketAnimationList, msg => {
            this.animations = msg.data.anims;
            console.log("Animations received:", this.animations);
        });
    }

    play(animation: Animation | null) {
        if (animation != null) {
            this.api.sendSelectMessage(animation.name);
            this.playingAnimation = animation;
        }

        return true;
    }

    stop() {
        this.api.sendSelectMessage("");
        this.playingAnimation = null;
    }
}