import { Car } from "./models/car.js";
import { Road } from "./models/road.js";
import { Visualizer } from "./models/visualizer.js";

const carCanvas = document.getElementById("carCanvas");
const networkCanvas = document.getElementById("networkCanvas");

carCanvas.width = 200;
networkCanvas.width = 500;


const carCtx = carCanvas.getContext("2d");
const networkCanvasCtx = networkCanvas.getContext("2d");

const road = new Road(carCanvas.width / 2, carCanvas.width, 3);
const N = 10;

const cars = generateCars(N);

const traffic = [
  new Car(road.getLaneCenter(1), -100, 30, 50, 'DUMMY', 2),
];


animate();

function generateCars(N) {
  const cars = [];
  for (let i = 1; i <= N; i++) {
    cars.push(new Car(road.getLaneCenter(1), 100, 30, 50, "AI"));
  }
  return cars;
}

function animate() {
  for (let i = 0; i < traffic.length; i++) {
    traffic[i].update(road.borders, []);
  }

  // PERF: setting carCanvas.height every frame forces full buffer reallocation
  // FIX: move to a window resize event listener instead
  carCanvas.height = window.innerHeight;
  networkCanvas.height = window.innerHeight;
  for (let i = 0; i < cars.length; i++) {
    cars[i].update(road.borders, traffic);
  }

  const bestCar = cars.find(car => car.y === Math.min(...cars.map(c => c.y)));

  carCtx.save();
  carCtx.translate(0, -bestCar.y + carCanvas.height * 0.7);

  road.draw(carCtx);
  for (let i = 0; i < traffic.length; i++) {
    traffic[i].draw(carCtx, 'red');
  }

  carCtx.globalAlpha = 0.2;
  for (let i = 0; i < cars.length; i++) {
    cars[i].draw(carCtx, 'blue');
  }
  carCtx.globalAlpha = 1;
  bestCar.draw(carCtx, 'blue', { drawSensor: true });

  carCtx.restore();

  Visualizer.drawNetwork(networkCanvasCtx, bestCar.brain);
  requestAnimationFrame(animate);
}
