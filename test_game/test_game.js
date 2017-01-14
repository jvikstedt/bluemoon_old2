var renderer = PIXI.autoDetectRenderer(800, 600,{backgroundColor : 0x1099bb});
document.body.appendChild(renderer.view);

var stage = new PIXI.Container();
var texture = PIXI.Texture.fromImage('http://pixijs.github.io/examples/required/assets/basics/bunny.png');

var bunny = new PIXI.Sprite(texture);

bunny.anchor.x = 0.5;
bunny.anchor.y = 0.5;

bunny.position.x = 200;
bunny.position.y = 150;

stage.addChild(bunny);

animate();
function animate() {
  requestAnimationFrame(animate);

  bunny.rotation += 0.1;

  renderer.render(stage);
}

var ws = new WebSocket("ws://localhost:4000/ws");
ws.onopen = function() {
  console.log("Connection open");
};
ws.onmessage = function (evt) {
  console.log("got message: " + evt.data);
};
ws.onclose = function() {
  console.log("Connection closed");
};
