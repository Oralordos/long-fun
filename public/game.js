var SCREEN_WIDTH = 800;
var SCREEN_HEIGHT = 600;

function Game() {
  PIXI.Container.call(this);
  this.units = [];
  this.map = null;
  this.stage = new PIXI.Container();
  this.renderer = PIXI.autoDetectRenderer(SCREEN_WIDTH, SCREEN_HEIGHT, {backgroundColor: 0x1099bb});
  document.querySelector('#game').appendChild(this.renderer.view);
  requestAnimationFrame(this.tick.bind(this));
}

Game.prototype = Object.create(PIXI.Container.prototype);
Game.prototype.constructor = Game;

Game.prototype.setState = function(gameState) {
  this.map = new Map(gameState.Map);
};

Game.prototype.tick = function() {
  requestAnimationFrame(this.tick.bind(this));
  this.renderer.render(this.stage);
};

function getGameState(callback) {
  var xhr = new XMLHttpRequest();
  xhr.open('GET', '/api/game/' + gameId);
  xhr.addEventListener('readystatechange', function() {
    if (xhr.readyState === 4) {
      if (xhr.status === 200) {
        callback(JSON.parse(xhr.responseText));
      } else {
        alert(xhr.responseText);
      }
    }
  });
  xhr.send();
}

function Map(mapData) {
  PIXI.Container.call(this);
  this.zoom = 1;
  this.interactive = true;
  this.on('mousedown', onDragStart)
      .on('touchstart', onDragStart)
      .on('mouseup', onDragEnd)
      .on('mouseupoutside', onDragEnd)
      .on('touchend', onDragEnd)
      .on('touchendoutside', onDragEnd)
      .on('mousemove', onDragMove)
      .on('touchmove', onDragMove);
  this.createMap(mapData);
  game.renderer.view.addEventListener('wheel', onScroll.bind(this));
  game.stage.addChild(this);
}

Map.prototype = Object.create(PIXI.Container.prototype);
Map.prototype.constructor = Map;

Map.prototype.createMap = function(mapData) {
  this.mapWidth = mapData.Width * mapData.Tileset.Tilewidth;
  this.mapHeight = mapData.Height * mapData.Tileset.Tileheight;
  this.setZoom(2, SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2);
  var self = this;

  loadTileset(mapData.Tileset, function() {
    for (var layer = 0; layer < mapData.Layers.length; layer++) {
      for (var x = 0; x < mapData.Width; x++) {
        for (var y = 0; y < mapData.Height; y++) {
          var tileID = mapData.Layers[layer][x + y * mapData.Width];
          if (tileID !== 0) {
            var sprite = PIXI.Sprite.fromFrame('tile-' + tileID);
            sprite.position.x = x * mapData.Tileset.Tilewidth;
            sprite.position.y = y * mapData.Tileset.Tileheight;
            self.addChild(sprite);
          }
        }
      }
    }
    this.cacheAsBitmap = true;
  });
};

Map.prototype.setZoom = function(s, mouseX, mouseY) {
  var maxX = this.mapWidth / SCREEN_WIDTH;
  var maxY = this.mapHeight / SCREEN_HEIGHT;
  var minZoom = Math.min(maxX, maxY);
  var oldZoom = this.zoom;
  var newZoom = Math.min(Math.max(s, minZoom), 6);
  this.zoom = newZoom;
  this.scale = new PIXI.Point(newZoom, newZoom);
  var zoomScale = newZoom / oldZoom;
  var newX = (this.position.x - mouseX) * zoomScale + mouseX;
  var newY = (this.position.y - mouseY) * zoomScale + mouseY;
  this.setPosition(newX, newY);
};

Map.prototype.setPosition = function(x, y) {
  this.position.x = Math.min(Math.max(x, SCREEN_WIDTH - this.mapWidth * this.zoom), 0);
  this.position.y = Math.min(Math.max(y, SCREEN_HEIGHT - this.mapHeight * this.zoom), 0);
};

function onTilesetLoaded(width, height, tileWidth, tileHeight, filename, callback) {
  var texture = PIXI.utils.TextureCache[filename];
  for (var y = 0; y < height; y++) {
    for (var x = 0; x < width; x++) {
      var frame = new PIXI.Rectangle(x * tileWidth, y * tileHeight, tileWidth, tileHeight);
      var tile = new PIXI.Texture(texture.baseTexture, frame);
      PIXI.Texture.addTextureToCache(tile, 'tile-' + (1 + x + y * width));
    }
  }
  callback();
}

function loadTileset(tileset, callback) {
  PIXI.loader.add('tiles', tileset.Filename)
             .once('complete', onTilesetLoaded.bind(null, tileset.Width, tileset.Height, tileset.Tilewidth, tileset.Tileheight, tileset.Filename, callback))
             .load();
}

function onScroll(e) {
  var amount;
  if (e.deltaY < 0) {
    amount = 1.25;
  } else {
    amount = 1 / 1.25;
  }
  this.setZoom(this.zoom * amount, e.clientX, e.clientY);
}

function onDragStart(e) {
  this.data = e.data;
  offsetWindow = this.data.getLocalPosition(this.parent);
  this.offsetX = offsetWindow.x - this.position.x;
  this.offsetY = offsetWindow.y - this.position.y;
  this.dragging = 1;
  this.originalTouch = offsetWindow;
}

function onDragEnd() {
  if (this.dragging === 1) {
    // TODO Activate unit
  }
  this.dragging = 0;
  this.data = null;
}

function onDragMove() {
  if (this.dragging === 1) {
    var newPos = this.data.getLocalPosition(this.parent);
    if (Math.abs(newPos.x - this.originalTouch.x) + Math.abs(newPos.y - this.originalTouch.y) > 5) {
      this.dragging = 2;
    }
  } else if (this.dragging === 2) {
    var newPosition = this.data.getLocalPosition(this.parent);
    this.setPosition(newPosition.x - this.offsetX, newPosition.y - this.offsetY);
  }
}

var game = new Game();
getGameState(game.setState);
