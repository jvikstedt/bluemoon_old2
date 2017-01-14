var renderer = PIXI.autoDetectRenderer(800, 600,{backgroundColor : 0x1099bb});
document.body.appendChild(renderer.view);

var stage = new PIXI.Container();
var container = new PIXI.Container();
var texture = PIXI.Texture.fromImage('http://pixijs.github.io/examples/required/assets/basics/bunny.png');

stage.addChild(container);

function animate() {
  requestAnimationFrame(animate);

  renderer.render(stage);
}

var players = {};
function newPlayer(id, x, y) {
  var bunny = new PIXI.Sprite(texture);

  bunny.anchor.x = 0.5;
  bunny.anchor.y = 0.5;

  bunny.position.x = x;
  bunny.position.y = y;

  players[id] = bunny;

  container.addChild(bunny);
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
      var pl = obj.payload;
      newPlayer(pl.id, pl.x, pl.y);
      break;
    case "move":
      var pl = obj.payload;
      move(pl.id, pl.x, pl.y);
      break;
    default:
      console.log("not found handler: " + obj.name);
  }
};
ws.onclose = function() {
  console.log("Connection closed");
};

function onKeyUp(key) {
  if (key.keyCode === 38) {
    ws.send(`{"name": "direction", "payload": {"axis": "y", "val": -1}}`);
  } else if (key.keyCode === 40) {
    ws.send(`{"name": "direction", "payload": {"axis": "y", "val": 1}}`);
  } else if (key.keyCode === 37) {
    ws.send(`{"name": "direction", "payload": {"axis": "x", "val": -1}}`);
  } else if (key.keyCode === 39) {
    ws.send(`{"name": "direction", "payload": {"axis": "x", "val": 1}}`);
  }
}

document.addEventListener('keyup', onKeyUp);
animate();
