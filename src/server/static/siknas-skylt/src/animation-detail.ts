import {autoinject} from 'aurelia-framework';
import {EventAggregator} from 'aurelia-event-aggregator';
import {WebsocketConnected, WebsocketDisconnected, AnimationSelectionChanged} from './messages';
import { Animation } from "./types";
import { AnimationListService } from "./animation-list-service";

@autoinject()
export class AnimationDetail {
  routeConfig;
  animation: Animation;

  constructor(private ea: EventAggregator, private service: AnimationListService) { }

  activate(params, routeConfig) {
    this.routeConfig = routeConfig;
    this.animation = this.service.getByName(params.name)
    // TODO: Send event that we selected the given name.
  }
}