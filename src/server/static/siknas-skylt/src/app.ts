import { Router, RouterConfiguration } from 'aurelia-router';
import { autoinject, observable } from 'aurelia-framework';
import { WSAPI } from './ws-api';

@autoinject()
export class App {
  router: Router;

  constructor(private api: WSAPI) { }

  blurValue: number = 10;
  @observable brightness: number = 128;

  /*
  configureRouter(config: RouterConfiguration, router: Router) {
    config.title = 'Animationer';
    config.map([
      { route: '', moduleId: 'no-selection', title: 'Select' },
      { route: 'animations/:name',  moduleId: 'animation-detail', name:'animations' }
    ]);

    this.router = router;
  }*/

  brightnessChanged(newValue, oldValue) {
    this.api.sendBrightnessMessage(parseInt(newValue));
  }

  created() {
    this.api.connect();
    // TODO: Add button to turn off animation.
  }
}
