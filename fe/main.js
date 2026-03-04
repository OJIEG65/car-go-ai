import { WSClient } from './ws-client.js';
import { Renderer } from './renderer.js';

const carCanvas = document.getElementById('carCanvas');
const networkCanvas = document.getElementById('networkCanvas');
carCanvas.width = 200;
networkCanvas.width = 500;

const renderer = new Renderer(carCanvas, networkCanvas);

// State from server
let road = null;
let worldState = null;
let brain = null;
let paused = false;

// Status elements
const statusEl = document.getElementById('status');
const pauseBtn = document.getElementById('pauseBtn');
const resetBtn = document.getElementById('resetBtn');
const saveBrainBtn = document.getElementById('saveBrainBtn');
const loadBrainBtn = document.getElementById('loadBrainBtn');

// Determine WebSocket URL from current page location
const wsProtocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
const wsUrl = `${wsProtocol}//${location.host}/ws`;

const client = new WSClient(wsUrl, {
  onConnect: () => {
    if (statusEl) statusEl.textContent = 'Connected';
    if (statusEl) statusEl.style.color = '#4f4';
  },

  onDisconnect: () => {
    if (statusEl) statusEl.textContent = 'Disconnected — reconnecting...';
    if (statusEl) statusEl.style.color = '#f44';
  },

  init: (payload) => {
    road = payload.road;
    // Fetch the best brain once connected
    fetchBrain();
  },

  state: (payload) => {
    worldState = payload;
  },

  generation_end: (payload) => {
    console.log(`Generation ${payload.generation} ended — best fitness: ${payload.bestFitness.toFixed(1)}`);
    // Refresh brain for network visualizer
    fetchBrain();
  },

  brain: (payload) => {
    brain = payload.data;
  },

  error: (payload) => {
    console.error('Server error:', payload.code, payload.message);
  },
});

// Fetch brain via REST for network visualization
async function fetchBrain() {
  try {
    const resp = await fetch('/api/brain/best');
    if (resp.ok) {
      brain = await resp.json();
    }
  } catch (e) {
    // Will retry on next generation
  }
}

// Periodically refresh brain for live visualization
setInterval(fetchBrain, 2000);

// Render loop
function animate() {
  renderer.drawFrame(worldState, road, brain);
  requestAnimationFrame(animate);
}
animate();

// Controls
if (pauseBtn) {
  pauseBtn.addEventListener('click', () => {
    paused = !paused;
    if (paused) {
      client.pause();
      pauseBtn.textContent = 'Resume';
    } else {
      client.resume();
      pauseBtn.textContent = 'Pause';
    }
  });
}

if (resetBtn) {
  resetBtn.addEventListener('click', () => {
    client.reset();
  });
}

if (saveBrainBtn) {
  saveBrainBtn.addEventListener('click', () => {
    client.saveBrain('best');
  });
}

if (loadBrainBtn) {
  loadBrainBtn.addEventListener('click', () => {
    client.loadBrain('best');
  });
}
