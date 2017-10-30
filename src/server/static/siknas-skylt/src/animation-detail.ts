import { autoinject } from 'aurelia-framework';
import { EventAggregator } from 'aurelia-event-aggregator';
import { WebsocketConnected, WebsocketDisconnected, AnimationViewed } from './messages';
import { Animation } from "./types";
import { AnimationListService } from "./animation-list-service";

@autoinject()
export class AnimationDetail {
    routeConfig;
    animation: Animation;

    constructor(private events: EventAggregator, private service: AnimationListService) { }

    activate(params, routeConfig) {
        this.routeConfig = routeConfig;
        this.animation = this.service.getByName(params.name)
        // TODO: Send event that we selected the given name.
        this.events.publish(new AnimationViewed(this.animation));
    }
}