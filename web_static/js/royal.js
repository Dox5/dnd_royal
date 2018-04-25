var config = {
    type: Phaser.AUTO,
    width: window.innerWidth,
    height: window.innerHeight,
    scene: [HostJoinScene, PlayingScene],
    backgroundColor: 0x222222
};

function resize(width, height) {
    if (width === undefined) { width = this.sys.game.config.width; }
    if (height === undefined) { height = this.sys.game.config.height; }

    if(this.cameras.main !== undefined) {
        this.cameras.main.setSize(width, height);
    }
}

var game = new Phaser.Game(config);

window.addEventListener('resize', function (event) {
    game.resize(window.innerWidth, window.innerHeight);
}, false);
