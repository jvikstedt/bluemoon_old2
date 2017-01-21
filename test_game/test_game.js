var renderer = PIXI.autoDetectRenderer(800, 600,{backgroundColor : 0x1099bb});
document.body.appendChild(renderer.view);

var stage = new PIXI.Container();
var container = new PIXI.Container();
var texture = PIXI.Texture.fromImage('http://pixijs.github.io/examples/required/assets/basics/bunny.png');
var appleTexture = PIXI.Texture.fromImage('https://cdn0.iconfinder.com/data/icons/glyph_set/128/apple.png');

stage.addChild(container);

function animate() {
  requestAnimationFrame(animate);

  renderer.render(stage);
}

var xDir = 0;
var yDir = 0;

var players = {};
var apples = {};
function newPlayer(id, x, y) {
  var bunny = new PIXI.Sprite(texture);

  bunny.anchor.x = 0.5;
  bunny.anchor.y = 0.5;

  bunny.position.x = x;
  bunny.position.y = y;

  players[id] = bunny;

  container.addChild(bunny);
}

function newApple(id, x, y) {
  var apple = new PIXI.Sprite(appleTexture);

  apple.anchor.x = 0.5;
  apple.anchor.y = 0.5;

  apple.width = 20;
  apple.height = 20;

  apple.position.x = x;
  apple.position.y = y;

  apples[id] = apple;

  container.addChild(apple);
}

function removePlayer(id) {
  container.removeChild(players[id]);
  delete players[id];
}

function move(id, x, y) {
  players[id].position.x = x;
  players[id].position.y = y;
}

var ws = new WebSocket("ws://localhost:4000/ws");
ws.onopen = function() {
  console.log("Connection open");
};
ws.onmessage = function (evt) {
  stage.removeChild(players[1]);
  var obj = JSON.parse(evt.data);
  console.log(obj);

  switch(obj.name) {
    case "new_player":
      newPlayer(obj.id, obj.x, obj.y);
      break;
    case "remove_player":
      removePlayer(obj.id);
    case "move":
      move(obj.id, obj.x, obj.y);
      break;
    case "new_apple":
      newApple(obj.id, obj.x, obj.y);
      break;
    default:
      console.log("not found handler: " + obj.name);
  }
};
ws.onclose = function() {
  console.log("Connection closed");
};

function onKeyUp(key) {
  if (key.keyCode === 38 && yDir == -1) {
    yDir = 0;
    ws.send(`{"name": "change_dir", "payload": {"axis": "y", "val": 0}}\n`);
  } else if (key.keyCode === 40 && yDir == 1) {
    yDir = 0;
    ws.send(`{"name": "change_dir", "payload": {"axis": "y", "val": 0}}\n`);
  } else if (key.keyCode === 37 && xDir == -1) {
    xDir = 0;
    ws.send(`{"name": "change_dir", "payload": {"axis": "x", "val": 0}}\n`);
  } else if (key.keyCode === 39 && xDir == 1) {
    xDir = 0;
    ws.send(`{"name": "change_dir", "payload": {"axis": "x", "val": 0}}\n`);
  }
}

function onKeyDown(key) {
  if (key.keyCode === 38 && yDir != -1) {
    yDir = -1;
    ws.send(`{"name": "change_dir", "payload": {"axis": "y", "val": -1}}\n`);
  } else if (key.keyCode === 40 && yDir != 1) {
    yDir = 1;
    ws.send(`{"name": "change_dir", "payload": {"axis": "y", "val": 1}}\n`);
  } else if (key.keyCode === 37 && xDir != -1) {
    xDir = -1;
    ws.send(`{"name": "change_dir", "payload": {"axis": "x", "val": -1}}\n`);
  } else if (key.keyCode === 39 && xDir != 1) {
    xDir = 1;
    ws.send(`{"name": "change_dir", "payload": {"axis": "x", "val": 1}}\n`);
  }
}

document.addEventListener('keyup', onKeyUp);
document.addEventListener('keydown', onKeyDown);
animate();
