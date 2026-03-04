import { lerp } from './utils.js';

// Draws the road, cars, sensors, and traffic from server state.
export class Renderer {
  constructor(carCanvas, networkCanvas) {
    this.carCanvas = carCanvas;
    this.networkCanvas = networkCanvas;
    this.carCtx = carCanvas.getContext('2d');
    this.netCtx = networkCanvas.getContext('2d');
  }

  resize() {
    this.carCanvas.height = window.innerHeight;
    this.networkCanvas.height = window.innerHeight;
  }

  drawFrame(state, road, brain) {
    if (!state || !road) return;

    this.resize();
    const ctx = this.carCtx;
    const best = state.cars[state.bestIndex];
    if (!best) return;

    ctx.save();
    ctx.translate(0, -best.y + this.carCanvas.height * 0.7);

    this.drawRoad(ctx, road);
    this.drawTraffic(ctx, state.traffic);
    this.drawCars(ctx, state.cars, state.bestIndex);
    this.drawSensors(ctx, state.bestSensor);

    ctx.restore();

    // Draw HUD
    this.drawHUD(ctx, state);

    // Draw neural network
    if (brain) {
      this.drawNetwork(this.netCtx, brain);
    }
  }

  drawRoad(ctx, road) {
    ctx.lineWidth = 5;
    ctx.strokeStyle = 'white';

    // Lane dividers
    for (let i = 1; i <= road.laneCount - 1; i++) {
      const x = lerp(road.left, road.right, i / road.laneCount);
      ctx.setLineDash([20, 20]);
      ctx.beginPath();
      ctx.moveTo(x, road.top);
      ctx.lineTo(x, road.bottom);
      ctx.stroke();
    }

    // Borders
    ctx.setLineDash([]);
    for (const border of road.borders) {
      ctx.beginPath();
      ctx.moveTo(border[0].x, border[0].y);
      ctx.lineTo(border[1].x, border[1].y);
      ctx.stroke();
    }
  }

  drawTraffic(ctx, traffic) {
    if (!traffic) return;
    for (const car of traffic) {
      this.drawCarPolygon(ctx, car, 'red');
    }
  }

  drawCars(ctx, cars, bestIndex) {
    if (!cars) return;
    ctx.globalAlpha = 0.2;
    for (let i = 0; i < cars.length; i++) {
      this.drawCarPolygon(ctx, cars[i], 'blue');
    }
    ctx.globalAlpha = 1;
    if (cars[bestIndex]) {
      this.drawCarPolygon(ctx, cars[bestIndex], 'blue');
    }
  }

  drawCarPolygon(ctx, car, color) {
    if (!car.polygon || car.polygon.length === 0) return;
    ctx.fillStyle = car.damaged ? 'grey' : color;
    ctx.beginPath();
    ctx.moveTo(car.polygon[0].x, car.polygon[0].y);
    for (let i = 1; i < car.polygon.length; i++) {
      ctx.lineTo(car.polygon[i].x, car.polygon[i].y);
    }
    ctx.fill();
  }

  drawSensors(ctx, sensorState) {
    if (!sensorState || !sensorState.rays) return;

    for (let i = 0; i < sensorState.rays.length; i++) {
      const ray = sensorState.rays[i];
      const reading = sensorState.readings ? sensorState.readings[i] : null;

      let end = ray[1];
      if (reading) {
        end = reading;
      }

      // Yellow line from start to hit/end
      ctx.beginPath();
      ctx.lineWidth = 2;
      ctx.strokeStyle = 'yellow';
      ctx.moveTo(ray[0].x, ray[0].y);
      ctx.lineTo(end.x, end.y);
      ctx.stroke();

      // Black line from ray end to hit point
      ctx.beginPath();
      ctx.lineWidth = 2;
      ctx.strokeStyle = 'black';
      ctx.moveTo(ray[1].x, ray[1].y);
      ctx.lineTo(end.x, end.y);
      ctx.stroke();
    }
  }

  drawHUD(ctx, state) {
    ctx.fillStyle = 'rgba(0,0,0,0.6)';
    ctx.fillRect(5, 5, 190, 50);
    ctx.fillStyle = 'white';
    ctx.font = '12px monospace';
    ctx.fillText(`Gen: ${state.generation}  Tick: ${state.tick}`, 12, 22);
    ctx.fillText(`Cars: ${state.cars.length}  Alive: ${state.cars.filter(c => !c.damaged).length}`, 12, 40);
  }

  // Neural network visualization (ported from visualizer.js)
  drawNetwork(ctx, network) {
    if (!network || !network.levels) return;
    const margin = 50;
    const left = margin;
    const top = margin;
    const width = ctx.canvas.width - margin * 2;
    const height = ctx.canvas.height - margin * 2;
    const levelHeight = height / network.levels.length;

    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

    for (let i = network.levels.length - 1; i >= 0; i--) {
      const levelTop = top + lerp(
        height - levelHeight, 0,
        network.levels.length === 1 ? 0.5 : i / (network.levels.length - 1)
      );
      ctx.setLineDash([7, 3]);
      const labels = i === network.levels.length - 1 ? ['\u2191', '\u2190', '\u2192', '\u2193'] : [];
      this.drawLevel(ctx, network.levels[i], left, levelTop, width, levelHeight, labels);
    }
  }

  drawLevel(ctx, level, left, top, width, height, outputLabels) {
    const right = left + width;
    const bottom = top + height;
    const { inputs, outputs, weights, biases } = level;

    // Weight connections
    for (let i = 0; i < inputs.length; i++) {
      for (let j = 0; j < outputs.length; j++) {
        ctx.beginPath();
        ctx.moveTo(this.getNodeX(inputs, i, left, right), bottom);
        ctx.lineTo(this.getNodeX(outputs, j, left, right), top);
        ctx.lineWidth = 2;
        ctx.strokeStyle = this.getRGBA(weights[i][j]);
        ctx.stroke();
      }
    }

    const nodeRadius = 18;

    // Input nodes
    for (let i = 0; i < inputs.length; i++) {
      const x = this.getNodeX(inputs, i, left, right);
      ctx.beginPath();
      ctx.arc(x, bottom, nodeRadius, 0, Math.PI * 2);
      ctx.fillStyle = 'black';
      ctx.fill();
      ctx.beginPath();
      ctx.arc(x, bottom, nodeRadius * 0.6, 0, Math.PI * 2);
      ctx.fillStyle = this.getRGBA(inputs[i]);
      ctx.fill();
    }

    // Output nodes
    for (let i = 0; i < outputs.length; i++) {
      const x = this.getNodeX(outputs, i, left, right);
      ctx.beginPath();
      ctx.arc(x, top, nodeRadius, 0, Math.PI * 2);
      ctx.fillStyle = 'black';
      ctx.fill();
      ctx.beginPath();
      ctx.arc(x, top, nodeRadius * 0.6, 0, Math.PI * 2);
      ctx.fillStyle = this.getRGBA(outputs[i]);
      ctx.fill();

      // Bias ring
      ctx.beginPath();
      ctx.lineWidth = 2;
      ctx.arc(x, top, nodeRadius * 0.8, 0, Math.PI * 2);
      ctx.strokeStyle = this.getRGBA(biases[i]);
      ctx.setLineDash([3, 3]);
      ctx.stroke();
      ctx.setLineDash([]);

      if (outputLabels[i]) {
        ctx.beginPath();
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillStyle = 'black';
        ctx.strokeStyle = 'white';
        ctx.font = (nodeRadius * 1.5) + 'px Arial';
        ctx.fillText(outputLabels[i], x, top + nodeRadius * 0.1);
        ctx.lineWidth = 0.5;
        ctx.strokeText(outputLabels[i], x, top + nodeRadius * 0.1);
      }
    }
  }

  getNodeX(nodes, index, left, right) {
    return lerp(left, right, nodes.length === 1 ? 0.5 : index / (nodes.length - 1));
  }

  getRGBA(value) {
    const alpha = Math.abs(value);
    const R = value < 0 ? 0 : 255;
    const G = R;
    const B = value > 0 ? 0 : 255;
    return `rgba(${R},${G},${B},${alpha})`;
  }
}
