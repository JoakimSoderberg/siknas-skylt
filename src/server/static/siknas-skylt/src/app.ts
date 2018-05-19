import { autoinject, observable } from 'aurelia-framework';
import { WSAPI } from './ws-api';

@autoinject()
export class App {
  constructor(private api: WSAPI) { }

  blurValue: number = 9;
  @observable brightness: number = 128;

  brightnessChanged(newValue, oldValue) {
    this.api.sendBrightnessMessage(parseInt(newValue));
  }

  created() {
    this.api.connect();
    // TODO: Add button to turn off animation.

    // TODO: Listen to server state here:
    /* events.subscribe(WebsocketAnimationList, msg_raw => {
            let msg: AnimationListMessage = msg_raw.data;
            console.log("Animations received:", msg); */
  }

  stop() {
    console.log("Stoppings");
    this.api.sendPlayMessage("");
  }
}
