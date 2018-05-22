import { autoinject, observable } from 'aurelia-framework';
import { WSAPI } from './ws-api';
import { EventAggregator } from 'aurelia-event-aggregator';
import { WebsocketAnimationList } from './messages';
import { AnimationListMessage } from './types';

@autoinject()
export class App {
  constructor(private api: WSAPI, private events: EventAggregator) { }

  blurValue: number = 9;
  @observable brightness: number = 128;

  isSending: boolean = false;

  brightnessChanged(newValue, oldValue) {
    this.isSending = true;
    this.api.sendBrightnessMessage(parseInt(newValue));
    this.isSending = false;
  }

  created() {
    this.api.connect();
    // TODO: Add button to turn off animation.

    this.events.subscribe(WebsocketAnimationList, msg_raw => {
      let msg: AnimationListMessage = msg_raw.data;
      console.log("Animations received by app:", msg);
      if (msg.brightness != this.brightness && !this.isSending) {
        this.brightness = msg.brightness;
      }
    });
  }

  stop() {
    console.log("Stopping");
    this.api.sendPlayMessage("");
  }
}
