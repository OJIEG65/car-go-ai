import { Car } from "./models/car.js";
import { Road } from "./models/road.js";

const canvas = document.getElementById("canvas");

canvas.width = 200;
canvas.height = window.innerHeight;
const ctx = canvas.getContext("2d");

const road = new Road(canvas.width / 2, canvas.width, 3);
const car = new Car(road.getLaneCenter(1), 100, 30, 50);

animate();

function animate() {
  canvas.height = window.innerHeight;
  car.update();

  ctx.save();
  ctx.translate(0, -car.y + canvas.height * 0.7);

  road.draw(ctx);
  car.draw(ctx);
  requestAnimationFrame(animate);
}
