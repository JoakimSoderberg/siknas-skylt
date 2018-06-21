

GUI
===

- [x] Add a live icon to indicate current animation
- [x] Add thumbnails to the list items
- [ ] Adjust to mobile screen vs big screen
- [x] Add control panel graphic
    - [ ] When control panel "owns" the server it becomes big and list fades out
    - [ ] Updates to show the status of the real control panel
    - [ ] Minimizes once "Custom" is selected and the webapp can be used
- [x] Host the thumbnails for the sketches
    - [x] Uses the same names as the animations in the config

Server
======

- [ ] Fix broadcaster with either mutex or channels
- [x] When control panel is connected
    - [x] Send control panel status to Websocket clients
    - [x] If the custom program is not selected disable choosing programs
        - [ ] Send command to websocket clients to grey out list
        - [x] Ignore websocket client selects server side

Thumbnail generator
===================

- [x] Start the thumbnail generating client
- [x] Make the thumbnail generating client switch sketches and record to SVG thumbnails and then quit
    - [x] Save thumbnails using same name as animations
