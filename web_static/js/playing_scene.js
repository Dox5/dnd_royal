var inverse_stencil = function(renderer, mask, camera) {
    var gl = renderer.gl;
    var geometryMask = this.geometryMask;

    // Force flushing before drawing to stencil buffer
    renderer.flush();

    // Enable and setup GL state to write to stencil buffer
    gl.enable(gl.STENCIL_TEST);
    gl.clear(gl.STENCIL_BUFFER_BIT);
    gl.colorMask(false, false, false, false);
    gl.stencilFunc(gl.EQUAL, 1, 1);
    gl.stencilOp(gl.REPLACE, gl.REPLACE, gl.REPLACE);

    // Write stencil buffer
    geometryMask.renderWebGL(renderer, geometryMask, 0.0, camera);
    renderer.flush();

    // Use stencil buffer to affect next rendering object
    gl.colorMask(true, true, true, true);
    gl.stencilFunc(gl.NOTEQUAL, 1, 1);
    gl.stencilOp(gl.INVERT, gl.INVERT, gl.INVERT);
}

function listenForDoubleClick(gameObject) {
    var alreadyClicked = false
    var doubleClickTimeout

    gameObject.on('pointerup', function(pointer) {
        if (alreadyClicked) {
            gameObject.emit("doubleclick")
            alreadyClicked = false
            clearTimeout(doubleClickTimeout)
        } else {
            alreadyClicked = true
            doubleClickTimeout = setTimeout(function() {
                alreadyClicked = false 
            }, 500)
        }
    })
}

var PlayerTokenController = new Phaser.Class({
    radius: 25,
    textures: [],
    updateRateHz: 2,
    timeTillUpdate: 0,
    tokens: {},

    Extends: Phaser.EventEmitter,
    
    initialize: function(scene, roomId, textures) {
        this.roomId = roomId
        this.updatePeriodMS = (1 / this.updateRateHz) * 1000,
        this.timeTillUpdate = 0
        this.textures = textures
        this.scene = scene

        Phaser.EventEmitter.call(this)
    },

    createGameTokenFor: function(token) {
        sprite = this.scene.add.sprite(token.Position.X,
                                       token.Position.Y,
                                       this.textures.pop())

        gameToken = {
            sprite: sprite,
            updatePosition: true,
            id: token.Id
        }

        this.tokens[token.Id]  = gameToken
        
        this.emit("create", gameToken)
    },

    updateGameToken: function(token) {
        gameToken = this.tokens[token.Id]
        if (!gameToken.updatePosition) {
            return
        }

        gameToken.sprite.x = token.Position.X
        gameToken.sprite.y = token.Position.Y
    },

    pollState: function() {
        do_http_get("api/v1/token/getTokens?RoomId=" + this.roomId)
            .then((tokenData) => {
                tokenData.forEach((token) => {
                    gameToken = this.tokens[token.Id]

                    if (gameToken === undefined) {
                        this.createGameTokenFor(token)
                    } else {
                        this.updateGameToken(token)
                    }
                })
            }, console.log)
    },

    update: function(time, timeDelta) {
        if (isNaN(timeDelta)) {
            timeDelta = 0
        }

        this.timeTillUpdate -= timeDelta

        if (this.timeTillUpdate <= 0) {
            this.timeTillUpdate = this.updatePeriodMS
            this.pollState()
        }
    }
})

var PlayingScene = new Phaser.Class({
    min_zoom: 0.3,

    Extends: Phaser.Scene,

    initialize: function() {
        Phaser.Scene.call(this, { key: 'PlayingScene' });
    },

    contols: null,
    updaters: [],

    preload: function() {
        this.load.setBaseURL('http://localhost:8000');

        this.load.image('token1', 'assets/token1.png')
        this.load.image('token2', 'assets/token2.png')
        this.load.image('token3', 'assets/token3.png')
        this.load.image('map',   'assets/map-update-3.png');
        this.load.image('grid',  'assets/grid.png');
        this.load.image('fog',   'assets/fog-pattern.png');
        this.load.image('circle-centre', 'assets/circle-centre.png');
        this.load.image('radius-setter', 'assets/radius-setter.png');

        this.load.spritesheet('buttons', 'assets/buttons.png', {frameWidth: 96,
                                                                frameHeight: 56})
    },

    init: function(data) {
        console.log(data)
        this.roomId = data.roomId
        this.userId = data.userId
    },

    create: function() {
        var map = this.add.sprite(0, 0, 'map').setOrigin(0,0)
        var grid = this.add.sprite(0, 0, 'grid').setOrigin(0, 0)
        grid.alpha = 0.3

        this.events.on('resize', resize, this)

        
        this.input.topOnly = true

        function setBounds(zoomLevel) {
            zoomLevel = 1
            marginX = map.displayWidth * 0.1
            marginY = map.displayHeight * 0.1

            camWidth = this.cameras.main.width
            camHeight = this.cameras.main.height

            this.cameras.main.setBounds((map.x - marginX - camWidth/2) * zoomLevel,
                                        (map.y - marginY - camHeight/2) * zoomLevel,
                                        (map.displayWidth + marginX*2 + camWidth) * zoomLevel,
                                        (map.displayHeight + marginY*2 + camHeight) * zoomLevel)
        }

        setBounds.call(this, 1)
        this.cameras.main.on("zoom", setBounds, this)
        

        this.cameras.main.setBackgroundColor("#b8e0d9")

        this.updaters.push(this.camera_controls([map, grid]))
        this.updaters.push(this.createTargetDisplay())
        this.updaters.push(this.createFog())

        this.playerTokenController = new PlayerTokenController(this,
                                                               this.roomId,
                                                               ['token1',
                                                               'token2',
                                                               'token3'])

        if(typeof(this.userId) !== 'undefined') {
            this.updaters.push(this.createTargetSetter())
            this.createGMControls()
        }

        this.updaters.push(this.playerTokenController.update.bind(this.playerTokenController))
    },

    endError: function(error) {
        console.log(error)    
        this.scene.start("HostJoinScene")
    },

    createGMControls: function() {
        var pauseResumeHandler = function(pauseOrResume) {
            var request = {
                RoomId: this.roomId,
                GameMasterId: this.userId
            }
            
            var stringRequest = JSON.stringify(request)

            if (pauseOrResume == "pause") {
                endpoint = "/pause"
            } else if (pauseOrResume == "resume") {
                endpoint = "/resume"
            }

            do_http_post("api/v1/fog" + endpoint, stringRequest)
                .then(undefined, console.log)
        }

        this.playerTokenController.on("create", function(gameToken) {
            gameToken.sprite.setInteractive()
            this.input.setDraggable(gameToken.sprite)
            gameToken.sprite.on("drag", makeDragHandler(gameToken.sprite), this)

            gameToken.sprite.on("dragstart", function() {
                gameToken.updatePosition = false
            }, this)

            gameToken.sprite.on("dragend", function() {
                gameToken.updatePosition = true

                request = {
                    TokenId: "" + gameToken.id,
                    RoomId: this.roomId,
                    GameMasterId: this.userId,
                    Position: {X: gameToken.sprite.x,
                               Y: gameToken.sprite.y}
                }

                do_http_post("api/v1/token/setTokenPosition",
                            JSON.stringify(request))
                    .then(undefined, console.log)
            }, this)
        }, this)

    
        pauseResumeHandler = pauseResumeHandler.bind(this)

        function getPeriodFromUser(helpText) {
            re = /([0-9]+(?:\.[0-9]+)?)([hms])/

            // Accept values like 11s or 2h or 5m or 4.5h
            match = getUserInput(helpText + "(eg 11s, 2h, 5m, 120s, 1.5h)", re)
            

            if (match === undefined) {
                return
            }

            var time = match[1]
            var unit = match[2]

            var unitMapping = {
                "h": 60 * 60,
                "m": 60,
                "s": 1
            }

            var periodSeconds = time * unitMapping[unit]
            
            return periodSeconds
        }

        // resume button
        var resume = makeButton(this, {idle: 0, click: 1}, {x: -10925, y: -975})
            .on("click", () => {pauseResumeHandler("resume")})

        // pause button
        var pause = makeButton(this, {idle: 2, click: 3}, {x: -10925, y: -915})
            .on("click", () => {pauseResumeHandler("pause")})

        // period button
        var period = makeButton(this, {idle: 4, click: 5}, {x: -10925, y: -855})
            .on("click", function() {
                time = getPeriodFromUser("Travel Time: ")

                var request  = {
                    Period: time,
                    RoomId: this.roomId,
                    GameMasterId: this.userId
                }

                do_http_post("api/v1/fog/setPeriod", JSON.stringify(request))
                    .then(undefined, console.log)
            }.bind(this))

        // advance button
        var advance = makeButton(this, {idle: 6, click: 7}, {x: -10925, y: -795})
            .on("click", function() {
                time = getPeriodFromUser("Advance by: ")

                var request = {
                    Amount: time,
                    RoomId: this.roomId,
                    GameMasterId: this.userId
                }

                do_http_post("api/v1/fog/advanceTime", JSON.stringify(request))
                    .then(undefined, console.log)
            }.bind(this))

        this.cameras.main.ignore([resume, pause, period, advance])

        var gmControlCam = this.cameras.add(0, 0, 150, 250)
        gmControlCam.scrollX -= 11000
        gmControlCam.scrollY -= 1000
        gmControlCam.transparent = false
        gmControlCam.setBackgroundColor("#222244")
    },

    createTargetDisplay: function() {
        var targetCircle = new Phaser.Geom.Circle(0, 0, 0)
        var gfx = this.make.graphics()
        var updatePeriodMs = 1000
        var timeToUpdate = 0

        return function(time, timeDelta) {
            if( isNaN(timeDelta) ) {
                timeDelta = 0
            }

            timeToUpdate -= timeDelta

            if (timeToUpdate <= 0) {
                var _this = this
                do_http_get("api/v1/fog/getTarget?RoomId="+this.roomId, {})
                    .then(function(response) {
                        targetCircle.x = response.Target.Centre.X
                        targetCircle.y = response.Target.Centre.Y
                        targetCircle.radius = response.Target.Radius
                    }, function(error) {
                        _this.endError(error)
                    })
                timeToUpdate = updatePeriodMs
            }

            gfx.clear()
            gfx.lineStyle(5, 0x222233, 0.4)
            gfx.strokeCircleShape(targetCircle)
        }
    },

    createTargetSetter: function() {
        var baseScale = 0.5
        var circleCentre = this.add.sprite(100, 100, "circle-centre")
        circleCentre.setScale(baseScale)
        var radiusSetter = this.add.sprite(20, 20, "radius-setter")
        radiusSetter.setScale(baseScale)
        var zoomFactor = 1

        function calc_radius() {
            x = circleCentre.x - radiusSetter.x
            y = circleCentre.y - radiusSetter.y

            return Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2))
        }

        var outline = new Phaser.Geom.Circle(circleCentre.x,
                                             circleCentre.y,
                                             calc_radius())

        function updateOutline() {
            outline.x = circleCentre.x
            outline.y = circleCentre.y
            outline.radius = calc_radius()
        }


        circleCentre.setInteractive()
        this.input.setDraggable(circleCentre)
        
        listenForDoubleClick(circleCentre)
        circleCentre.on('doubleclick', function() {
            newTarget = {
                RoomId: this.roomId,
                GameMasterId: this.userId,
                FogTarget: {
                    Centre: {
                        X: circleCentre.x,
                        Y: circleCentre.y,
                    },
                    Radius: calc_radius()
                }
            }
            do_http_post("api/v1/fog/setTarget", JSON.stringify(newTarget))
        }, this)

        circleCentre.on('drag',
                        makeDragHandler(circleCentre,
                                        [radiusSetter],
                                        updateOutline),
                        this)

        radiusSetter.setInteractive()
        this.input.setDraggable(radiusSetter)
        radiusSetter.on('drag',
                        makeDragHandler(radiusSetter, [], updateOutline),
                        this)


        this.input.dragTimeThreshold = 100 // MS

        var gfx = this.make.graphics() 
        return function(time, timeDelta) {
            gfx.clear()
            gfx.lineStyle(2, 0x663333, 1)
            gfx.strokeCircleShape(outline)

            // Keep the size of the controls the same, always
            zoom = 1/this.cameras.main.zoom
            circleCentre.setScale(baseScale * zoom)
            radiusSetter.setScale(baseScale * zoom)
        }
    },

    createFog: function() {
        // Workout the largest possible area the payer can view an use that to
        // tile the fog
        function setFogSize() {
            bounds = this.cameras.main._bounds
            cfg = this.scene.systems.game.config
            halfViewableWidth = cfg.width / (this.min_zoom * 2)
            halfViewableHeight = cfg.height / (this.min_zoom * 2)

            fogDim = {
                left:   bounds.left - halfViewableWidth,
                top:    bounds.top - halfViewableHeight,
                right:  bounds.right + halfViewableWidth,
                bottom: bounds.bottom + halfViewableHeight
            }

            width = fogDim.right - fogDim.left
            height = fogDim.bottom - fogDim.top

            fog.setPosition(fogDim.left, fogDim.top)
            fog.setSize(width, height)
        }

        var fog = this.add.tileSprite(0, 0, 1024, 1024, 'fog')
                    .setOrigin(0, 0)

        this.events.on("resize", setFogSize, this)
        // Initialise the fog size
        setFogSize.call(this)

        

        var maskShape = new Phaser.Geom.Circle(window.innerWidth/2,
                                          window.innerHeight/2,
                                          300)

        var maskGfx = this.make.graphics()
        var rimGfx = this.make.graphics()

        var rate = {position: {x: 0, y: 0}, radius: 0}

        fog.mask = new Phaser.Display.Masks.GeometryMask(this, maskGfx)
        // This monkey patch is so hacky, I'm probably going to get annoyed by this
        // breaking wierdly one day...
        fog.mask.preRenderWebGL = inverse_stencil

        fog.alpha =1

        timeTillNextPollMs = 0

        return function(time, timeDelta) {
            pollPeriodMs = 2000
            timeTillNextPollMs -= timeDelta

            if(timeTillNextPollMs < 0) {
                timeTillNextPollMs = pollPeriodMs
                do_http_get("api/v1/fog/location?RoomId="+this.roomId, null).then(
                    function(loc) {
                        rate.radius = loc.Rate.Radius
                        rate.position.x = loc.Rate.Translation.X
                        rate.position.y = loc.Rate.Translation.Y

                        // Fix the position of the circle (It's likely going to be
                        // out!)
                        maskShape.x = loc.Current.Centre.X
                        maskShape.y = loc.Current.Centre.Y
                        maskShape.radius = loc.Current.Radius
                    }.bind(this), console.log
                )
            }
            
            // TODO: The resetting of the position upsets this because it is out
            // of sync
            //timeDeltaSeconds = timeDelta/1000
            //maskShape.x += rate.position.x * timeDeltaSeconds
            //maskShape.y += rate.position.y * timeDeltaSeconds
            //maskShape.radius += rate.radius * timeDeltaSeconds

            maskGfx.clear()
            maskGfx.fillStyle(0xffffff, 0)
            maskGfx.fillCircleShape(maskShape)

            rimGfx.clear()
            var rimWidth = 14
            rimGfx.lineStyle(rimWidth, 0xd048c8, 0.6)
            rimGfx.strokeCircle(maskShape.x,
                                maskShape.y,
                                maskShape.radius - rimWidth/2)
        }
    },
    
    timeTillNextPollMs: 0, 
    update: function(time, timeDelta) {


        this.updaters.forEach(updater => updater.call(this, time, timeDelta))
    },

    camera_controls: function(backgroundObjects) {
        var mouseDown = false
        var dragStart
        var dragEnd
        var zoom_step = 0.05
        var _this = this

        backgroundObjects.forEach(function(gameObject) {
            gameObject.setInteractive()


            gameObject.on('pointerdown', function(pointer, gameObject) {
                mouseDown = true
                dragStart = this.cameras.main.getWorldPoint(pointer.x,
                                                            pointer.y)
            }, _this)

            gameObject.on('pointerup', function(pointer, gameObject) {
                mouseDown = false
                dragStart = null
            }, _this)

            gameObject.on('pointermove', function(pointer) {
                if (mouseDown) {
                    dragHere = this.cameras.main.getWorldPoint(pointer.x,
                                                               pointer.y)

                    this.cameras.main.scrollX += dragStart.x - dragHere.x
                    this.cameras.main.scrollY += dragStart.y - dragHere.y

                    dragStart = this.cameras.main.getWorldPoint(pointer.x,
                                                                pointer.y)
                }
            }, _this)
        })

        wheelHandler = function(wheel) {
            var goingUp = wheel.deltaY < 0

            this.cameras.main.zoom += zoom_step *( goingUp? 1 : - 1);
            this.cameras.main.zoom = Math.max(this.min_zoom, this.cameras.main.zoom)

            this.cameras.main.emit("zoom", this.cameras.main.zoom)
        }

        window.addEventListener('wheel', function(wheel) {
            wheelHandler.call(_this, wheel)
        })

        var keyInputs = this.input.keyboard.addKeys({
            'up':      Phaser.Input.Keyboard.KeyCodes.W,
            'down':    Phaser.Input.Keyboard.KeyCodes.S,
            'left':    Phaser.Input.Keyboard.KeyCodes.A,
            'right':   Phaser.Input.Keyboard.KeyCodes.D,
            'zoomIn':  Phaser.Input.Keyboard.KeyCodes.Q,
            'zoomOut': Phaser.Input.Keyboard.KeyCodes.E
        });

        var controlConfig = {
            camera: this.cameras.main,
            left: keyInputs.left,
            right: keyInputs.right,
            up: keyInputs.up,
            down: keyInputs.down,
            zoomIn: keyInputs.zoomIn,
            zoomOut: keyInputs.zoomOut,
            acceleration: 0.5,
            drag: 0.01,
            maxSpeed: 2.0
        };

        var controls = new Phaser.Cameras.Controls.Smoothed(controlConfig);

        return function(time, timeDelta) {
            controls.update(timeDelta)
        }
    }
});
