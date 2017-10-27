import { Router, RouterConfiguration } from 'aurelia-router';
import { inject } from 'aurelia-framework';
import { WSAPI } from './ws-api';

@inject(WSAPI)
export class App {
  router: Router;

  food = [
    { id: 0, name: 'Pizza' },
    { id: 1, name: 'Cake' },
    { id: 2, name: 'Steak' },
    { id: 3, name: 'Pasta' },
    { id: 4, name: 'Fries' }
  ];
  selectedMeal = null;

  constructor(public api: WSAPI) { }

  configureRouter(config: RouterConfiguration, router: Router) {
    config.title = 'Animationer';
    config.map([
      { route: '', moduleId: 'no-selection', title: 'Select' },
      { route: 'animations/:name',  moduleId: 'animation-detail', name:'animations' }
    ]);

    this.router = router;
  }
}
