import { autoinject } from 'aurelia-framework';
import { EventAggregator } from 'aurelia-event-aggregator';
import { Animation, AnimationListMessage } from "./types";
import { WSAPI } from "./ws-api";
import { WebsocketAnimationList, WebsocketConnected, WebsocketDisconnected } from "./messages";

@autoinject()
export class AnimationListService {
    animations; // TODO: Let this have a type.

    constructor(private events: EventAggregator, private api: WSAPI) {
        // Once the websocket is connected we'll receive a list of animations.
        events.subscribe(WebsocketAnimationList, msg => {
            this.animations = msg.data.anims;
            console.log("Animations service received:", this.animations);
        });

        // TODO: Remove test code...
        /*let jsonData: string = `
        {
            "message_type": "list",
            "anims": [
                {"name": "Arne", "description": "Arne weise"},
                {"name": "Blarg", "description": "Blarg blorg"}
            ]
        }`;

        let data = JSON.parse(jsonData);
        let animations: AnimationListMessage = data;
        this.events.publish(new WebsocketAnimationList(animations));*/
    }

    getByName(name: string): Animation {
        if (!this.animations)
            return null;

        for (let animation of this.animations) {
            if (animation.name == name)
                return animation;
        }
        return null;
    }

    setSelectedAnimation(name: string) {
        this.api.sendSelectMessage(name);
    }

    get isConnected(): boolean {
        return this.api.isConnected;
    }
}
