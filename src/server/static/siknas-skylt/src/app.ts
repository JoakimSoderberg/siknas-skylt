import { autoinject, observable } from 'aurelia-framework';
import { WSAPI } from './ws-api';
import { EventAggregator } from 'aurelia-event-aggregator';
import { WebsocketAnimationList, WebsocketBrightnessMessage } from './messages';
import { AnimationListMessage, BrightnessMessage } from './types';
import { AnimationList } from './animation-list';

@autoinject()
export class App {
  constructor(private api: WSAPI, private events: EventAggregator) { }

  // TODO: If no animation show error about not being configured correctly.

  blurValue: number = 9;
  @observable brightness: number = 128;

  enableSending: boolean = true;

  brightnessChanged(newValue, oldValue) {
    if (this.enableSending) {
      this.api.sendBrightnessMessage(parseInt(newValue));
    }
  }

  created() {
    this.api.connect();
    // TODO: Add button to turn off animation.
    this.events.subscribeOnce(WebsocketAnimationList, msg_raw => {
      let msg: AnimationListMessage = msg_raw.data;
      this.brightness = msg.brightness;
    });

    this.events.subscribe(WebsocketBrightnessMessage, msg_raw => {
      let msg: BrightnessMessage = msg_raw.data;
      console.log("Brightness received: ", this.brightness);

      // Make sure we don't resend any incoming brightness changes
      // because of the view binding of brightness by waiting a while
      // after each incoming message.
      this.enableSending = false;
      setTimeout(() => {
        this.enableSending = true;
      }, 500);

      if (msg.brightness != this.brightness) {
        this.brightness = msg.brightness;
      }
    });
  }

  stop() {
    console.log("Stopping");
    this.api.sendPlayMessage("");
  }
}
