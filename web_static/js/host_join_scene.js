function getRoom() {
    re = /[0-9]+/
    match = getUserInput("RoomId: ", re)

    console.log(match)
    if(match === undefined) {
        return undefined
    } else {
        return match[0]
    }
}

function getUser() {
    re = /[0-9]+/
    match = getUserInput("UserId: ", re)

    if(match === undefined) {
        return undefined
    } else {
        return match[0]
    }
}

var HostJoinScene = new Phaser.Class({
    Extends: Phaser.Scene,

    initialize: function() {
        Phaser.Scene.call(this, { key: 'HostJoinScene' });
    },

    preload: function() {
        this.load.spritesheet('buttons', 'assets/buttons.png', {frameWidth: 96,
                                                                frameHeight: 56})
    },

    create: function() {

        // Join room
        makeButton(this, {idle: 8, click: 9}, {x: 100, y: 100})
            .on("click", function() {
                var roomId = getRoom()
                do_http_get("api/v1/fog/location?RoomId="+roomId)
                    .then(function() {
                              this.scene.start("PlayingScene", {"roomId": roomId})
                          }.bind(this),
                          function() {
                              console.log("No room with key ", roomId, " was found")
                          })
            }, this)

        // Create Room
        makeButton(this, {idle: 10, click: 11}, {x: 100, y: 164})
            .on("click", function() {
                do_http_post("api/v1/room/create", {})
                    .then(function(createdRoom) {
                              this.scene.start("PlayingScene",
                                               {roomId: createdRoom.RoomId,
                                                userId: createdRoom.GameMasterId})
                          }.bind(this),
                          function(failureDetails) {
                              console.log("Failed to create room:",
                                          failureDetails)
                    })
            }, this)

        // Join as GM
        makeButton(this, {idle: 12, click: 13}, {x: 100, y: 228})
            .on("click", function() {
                var roomId = getRoom()
                var userId = getUser()
                do_http_get("api/v1/fog/location?RoomId="+roomId)
                    .then(function() {
                              this.scene.start("PlayingScene",
                                               {"roomId": roomId,
                                                "userId": userId})
                          }.bind(this),
                          function() {
                              console.log("No room with key ", roomId, " was found")
                          })
            }, this)
    }
})
