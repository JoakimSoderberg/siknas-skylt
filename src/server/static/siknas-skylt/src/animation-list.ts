import { EventAggregator } from 'aurelia-event-aggregator';
import { WSAPI } from './ws-api';
import {
    WebsocketConnected, WebsocketDisconnected, WebsocketMessageReceived,
    WebsocketError, WebsocketAnimationList
} from './messages';
import { inject } from 'aurelia-framework';
import { Animation, AnimationListMessage } from "./types";

@inject(WSAPI, EventAggregator)
export class AnimationList {
    animations;
    selectedName: string = null;
    isActive: boolean = false;
    ea: EventAggregator;

    constructor(api: WSAPI, ea: EventAggregator) {
        this.ea = ea;

        ea.subscribe(WebsocketConnected, () => this.isActive = true);
        ea.subscribe(WebsocketDisconnected, msg => this.isActive = false);

        // Once the websocket is connected we'll receive a list of animations.
        ea.subscribe(WebsocketAnimationList, msg => {
            this.animations = msg.data.anims;
            console.log("Animations sub:", this.animations);
        });
    }

    created() {
        // Connect the websocket.
        //this.api.connect();
        let jsonData: string = `
        {
            "message_type": "list",
            "anims": [
                {"name": "Arne", "description": "Arne weise"},
                {"name": "Arne", "description": "Arne weise"}
            ]
        }`;

        let data = JSON.parse(jsonData);

        let animations: AnimationListMessage = data;
        //animations.anims = data["anims"] as Array<Animation>;

        console.log("before anims:", animations);
        this.ea.publish(new WebsocketAnimationList(animations));
    }

    select(animation: Animation | null) {
        if (animation != null) {
            this.selectedName = animation.name;
        } else {
            this.selectedName = null;
        }
        return true;
    }
}