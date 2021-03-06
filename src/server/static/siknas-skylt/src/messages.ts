
import { Animation, AnimationListMessage } from "./types";


export class WebsocketConnected {
    constructor() { }
}

export class WebsocketDisconnected {
    constructor(public event: CloseEvent) { }
}

export class WebsocketError {
    constructor(public event: ErrorEvent) { }
}

export class WebsocketMessageReceived {
    constructor(public event: MessageEvent) { }
}

export class WebsocketAnimationList {
    constructor(public data) { }
}

export class AnimationViewed {
    constructor(public animation: Animation) { }
}

export class WebsocketOPCMessage {
    constructor(public data) { }
}

// TODO: Make data have an inteface
export class WebsocketControlPanelMessage {
    constructor(public data) {}
}

export class WebsocketBrightnessMessage {
    constructor(public data) {}
}
