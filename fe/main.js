import { Car } from "./models/car.js";
import { Road } from "./models/road.js";

const carCanvas = document.getElementById("carCanvas");
const networkCanvas = document.getElementById("networkCanvas");

carCanvas.width = 200;
networkCanvas.width = 300;


const carCtx = carCanvas.getContext("2d");
const networkCanvasCtx = networkCanvas.getContext("2d");

const road = new Road(carCanvas.width / 2, carCanvas.width, 3);
const car = new Car(road.getLaneCenter(1), 100, 30, 50, "AI");

const traffic = [
  new Car(road.getLaneCenter(1), -100, 30, 50, 'DUMMY', 2),
];


animate();

function animate() {
  for (let i = 0; i < traffic.length; i++) {
    traffic[i].update(road.borders, []);
  }

  // PERF: setting carCanvas.height every frame forces full buffer reallocation
  // FIX: move to a window resize event listener instead
  carCanvas.height = window.innerHeight;
  networkCanvas.height = window.innerHeight;
  car.update(road.borders, traffic);

  carCtx.save();
  carCtx.translate(0, -car.y + carCanvas.height * 0.7);

  road.draw(carCtx);
  for (let i = 0; i < traffic.length; i++) {
    traffic[i].draw(carCtx, 'red');
  }
  car.draw(carCtx, 'blue');

  carCtx.restore();
  requestAnimationFrame(animate);
}
