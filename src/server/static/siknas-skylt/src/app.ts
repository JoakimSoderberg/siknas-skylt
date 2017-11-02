import { Router, RouterConfiguration } from 'aurelia-router';
import { autoinject } from 'aurelia-framework';
import { WSAPI } from './ws-api';

@autoinject()
export class App {
  router: Router;

  constructor(private api: WSAPI) { }

  configureRouter(config: RouterConfiguration, router: Router) {
    config.title = 'Animationer';
    config.map([
      { route: '', moduleId: 'no-selection', title: 'Select' },
      { route: 'animations/:name',  moduleId: 'animation-detail', name:'animations' }
    ]);

    this.router = router;
  }

  created() {
    this.api.connect();
  }
}
