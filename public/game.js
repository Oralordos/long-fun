var SCREEN_WIDTH = 800;
var SCREEN_HEIGHT = 600;

var gameMap;

function Map(mapWidth, mapHeight) {
  PIXI.Container.call(this);
  this.tiles = new PIXI.Container();
  this.addChild(this.tiles);
  this.mapWidth = mapWidth;
  this.mapHeight = mapHeight;
  this.zoom = 1;
  this.cacheAsBitmap = true;
  this.interactive = true;
  this.on('mousedown', onDragStart)
           .on('touchstart', onDragStart)
           .on('mouseup', onDragEnd)
           .on('mouseupoutside', onDragEnd)
           .on('touchend', onDragEnd)
           .on('touchendoutside', onDragEnd)
           .on('mousemove', onDragMove)
           .on('touchmove', onDragMove);
}

Map.prototype = PIXI.Container.prototype;
Map.prototype.Contructor = Map;

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

function animate() {
  requestAnimationFrame(animate);
  renderer.render(stage);
}

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

function loadTileset(width, height, tileWidth, tileHeight, filename, callback) {
  PIXI.loader.add('tiles', filename)
             .once('complete', onTilesetLoaded.bind(null, width, height, tileWidth, tileHeight, filename, callback))
             .load();
}

function loadMap(callback) {
  var xhr = new XMLHttpRequest();
  xhr.open('GET', '/maps/' + map);
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

function createMap(mapData) {
  var mapWidth = mapData.Tilewidth * mapData.Width;
  var mapHeight = mapData.Tileheight * mapData.Height;

  var tilesX = mapData.Tilesets[0].Imagewidth / mapData.Tilewidth;
  var tilesY = mapData.Tilesets[0].Imageheight / mapData.Tileheight;

  gameMap = new Map(mapWidth, mapHeight);

  loadTileset(tilesX, tilesY, mapData.Tilewidth, mapData.Tileheight, '/public/images/toens-medieval-strategy.png', function() {
    for (var layer = 0; layer < mapData.Layers.length; layer++) {
      if (mapData.Layers[layer].Name != 'UnitLayer') {
        for (var x = 0; x < mapData.Width; x++) {
          for (var y = 0; y < mapData.Height; y++) {
            var tileId = mapData.Layers[layer].Data[x + y * mapData.Width];
            if (tileId !== 0) {
              var sprite = PIXI.Sprite.fromFrame('tile-' + tileId);
              // TODO Setup animation
              sprite.position.x = x * mapData.Tilewidth;
              sprite.position.y = y * mapData.Tileheight;
              gameMap.tiles.addChild(sprite);
            }
          }
        }
      } else {
        // TODO Load units
      }
    }
    gameMap.setZoom(2, SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2);
    document.querySelector('canvas').addEventListener('wheel', onScroll.bind(gameMap));
    stage.addChild(gameMap);
  });
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

var stage = new PIXI.Container();
var renderer = PIXI.autoDetectRenderer(SCREEN_WIDTH, SCREEN_HEIGHT);
document.querySelector('#game').appendChild(renderer.view);
loadMap(createMap);
requestAnimationFrame(animate);
