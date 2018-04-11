

GUI
===

- [ ] Add a live icon to indicate current animation
- [ ] Add thumbnails to the list items
- [ ] Adjust to mobile screen vs big screen
- [ ] Add control panel graphic
    - [ ] When control panel "owns" the server it becomes big and list fades out
    - [ ] Updates to show the status of the real control panel
    - [ ] Minimizes once "Custom" is selected and the webapp can be used
- [ ] Host the thumbnails for the sketches
    - [ ] Uses the same names as the animations in the config

Server
======

- [ ] Fix broadcaster with either mutex or channels
- [ ] When control panel is connected
    - [ ] Send control panel status to Websocket clients
    - [ ] If the custom program is not selected disable choosing programs
        - [ ] Send command to websocket clients to grey out list
        - [ ] Ignore websocket client selects server side

Thumbnail generator
===================

- [ ] Start the thumbnail generating client
- [ ] Make the thumbnail generating client switch sketches and record to SVG thumbnails and then quit
    - [ ] Save thumbnails using same name as animations