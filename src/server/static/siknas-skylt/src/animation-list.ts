import { EventAggregator } from 'aurelia-event-aggregator';
import {
    WebsocketConnected, WebsocketDisconnected, WebsocketMessageReceived,
    WebsocketError, WebsocketAnimationList, AnimationSelectionChanged
} from './messages';
import { autoinject } from 'aurelia-framework';
import { Animation, AnimationListMessage } from "./types";
import { AnimationListService } from "./animation-list-service";

@autoinject()
export class AnimationList {
    selectedName: string = null;

    constructor(private events: EventAggregator, private service: AnimationListService) {

    }

    get animations() {
        return this.service.animations;
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