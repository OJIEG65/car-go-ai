import { Car } from "./models/car.js";
import { Road } from "./models/road.js";

const canvas = document.getElementById("canvas");

canvas.width = 200;
canvas.height = window.innerHeight;
const ctx = canvas.getContext("2d");

const road = new Road(canvas.width / 2, canvas.width, 3);
const car = new Car(road.getLaneCenter(1), 100, 30, 50, "AI");

const traffic = [
  new Car(road.getLaneCenter(1), -100, 30, 50, 'DUMMY', 2),
];


animate();

function animate() {
  for (let i = 0; i < traffic.length; i++) {
    traffic[i].update(road.borders, []);
  }

  // PERF: setting canvas.height every frame forces full buffer reallocation
  // FIX: move to a window resize event listener instead
  canvas.height = window.innerHeight;
  car.update(road.borders, traffic);

  ctx.save();
  ctx.translate(0, -car.y + canvas.height * 0.7);

  road.draw(ctx);
  for (let i = 0; i < traffic.length; i++) {
    traffic[i].draw(ctx, 'red');
  }
  car.draw(ctx, 'blue');

  ctx.restore();
  requestAnimationFrame(animate);
}
