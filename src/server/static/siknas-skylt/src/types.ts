export interface Animation {
    name: string
    description: string
}

export interface AnimationListMessage {
    playing: number
    playing_name: string
    // TODO: Add brightness
    anims: Array<Animation>
}

export interface ControlPanelMessage {
    program: number
    color: number[]
    brightness: number
}
