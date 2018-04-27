export interface Animation {
    name: string
    description: string
}

export interface AnimationListMessage {
    anims: Array<Animation>
}

export interface ControlPanelMessage {
    program: number
    color: number[]
    brightness: number
}
