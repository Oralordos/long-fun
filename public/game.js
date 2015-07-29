var SCREEN_WIDTH = 800;
var SCREEN_HEIGHT = 600;

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
  var container = new PIXI.Container();
  container.mapWidth = mapData.Tilewidth * mapData.Width;
  container.mapHeight = mapData.Tileheight * mapData.Height;
  var tilesX = mapData.Tilesets[0].Imagewidth / mapData.Tilewidth;
  var tilesY = mapData.Tilesets[0].Imageheight / mapData.Tileheight;
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
              container.addChild(sprite);
            }
          }
        }
      }
    }
    container.cacheAsBitmap = true;
    container.interactive = true;
    container.on('mousedown', onDragStart)
             .on('touchstart', onDragStart)
             .on('mouseup', onDragEnd)
             .on('mouseupoutside', onDragEnd)
             .on('touchend', onDragEnd)
             .on('touchendoutside', onDragEnd)
             .on('mousemove', onDragMove)
             .on('touchmove', onDragMove);
    stage.addChild(container);
  });
}

function onDragStart(e) {
  this.data = e.data;
  offsetWindow = this.data.getLocalPosition(this.parent);
  this.offsetX = offsetWindow.x - this.position.x;
  this.offsetY = offsetWindow.y - this.position.y;
  this.dragging = true;
}

function onDragEnd() {
  this.dragging = false;
  this.data = null;
}

function onDragMove() {
  if (this.dragging) {
    var newPosition = this.data.getLocalPosition(this.parent);
    this.position.x = Math.min(Math.max(newPosition.x - this.offsetX, SCREEN_WIDTH - this.mapWidth), 0);
    this.position.y = Math.min(Math.max(newPosition.y - this.offsetY, SCREEN_HEIGHT - this.mapHeight), 0);
  }
}

var stage = new PIXI.Container();
var renderer = PIXI.autoDetectRenderer(SCREEN_WIDTH, SCREEN_HEIGHT);
document.querySelector('#game').appendChild(renderer.view);
loadMap(createMap);
requestAnimationFrame(animate);
