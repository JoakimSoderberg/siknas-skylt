import {inject} from 'aurelia-framework';
import {EventAggregator} from 'aurelia-event-aggregator';
import {WSAPI} from './ws-api';
import {WebsocketConnected, WebsocketDisconnected, AnimationSelectionChanged} from './messages';
import { Animation } from "./types";
import { AnimationList } from "./animation-list";

interface Contact {
  firstName: string;
  lastName: string;
  email: string;
}

@inject(WSAPI, EventAggregator, AnimationList)
export class AnimationDetail {
  routeConfig;
  animation: Animation;

  constructor(private api: WSAPI, private ea: EventAggregator) { }

  activate(params, routeConfig) {
    this.routeConfig = routeConfig;
    this.ea.subscribe(AnimationSelectionChanged, msg => {
        console.log("Animation selection changed:", msg);
        this.animation = msg.data.animation;
    });
  }
}