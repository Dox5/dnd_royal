function do_http_request(uri, method, data) {
    return new Promise((resolve, reject) => {
        var request = new XMLHttpRequest();

        request.onreadystatechange = function() {
            if (request.readyState == 4) {
                if(request.status == 200) {
                    resolve(request.response)
                } else {
                    reject({status: request.status,
                            statusText: request.statusText})
                }
            }
        }
        
        request.responseType = "json"
        request.open(method, uri, true)
        request.send(data)
    })
}

function do_http_post(uri, data) {
    return do_http_request(uri, "POST", data)
}

function do_http_get(uri, data) {
    return do_http_request(uri, "GET", data)
}

function enableClick(gameObject, context) {
        gameObject.setInteractive()
        gameObject.on("pointerup", function() {
            gameObject.emit("click")
        })
}

function makeButton(game, frames, position) {
    var btn = game.add.sprite(position.x, position.y, 'buttons')
    enableClick(btn)

    btn.setFrame(frames.idle)
    btn.on("pointerdown", function() {
        btn.setFrame(frames.click)

        btn.once("click", function() {
            this.setFrame(frames.idle) 
        }, btn)
        
        btn.once("pointerout", function() {
            this.setFrame(frames.idle)
        }, btn)
    }, btn)

    return btn
}

function makeDragHandler(obj, children, extraCallback) {
    return function(pointer) {
        pointerWorld = this.cameras.main.getWorldPoint(pointer.x,
                                                       pointer.y)
        drag = {x: pointerWorld.x - obj.x,
                y: pointerWorld.y - obj.y}

        obj.x += drag.x
        obj.y += drag.y

        if (children === undefined) { children = [] }
        children.forEach(function(child) {
            child.x += drag.x 
            child.y += drag.y
        })

        if (extraCallback !== undefined) {
            extraCallback()
        }
    }
}

function getUserInput(promptText, matches) {
    while (true) {
        value = prompt(promptText)
        if (value === null) {
            // Cancel pressed
            return
        }

        if (matches === undefined) {
            console.log("matches:", matches)
            return value
        }

        match = value.match(re)

        if (match) {
            break;
        } else {
            console.log("Invalid!")
        }
    }
    return match
}
