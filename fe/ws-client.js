// WebSocket client with auto-reconnect.
export class WSClient {
  constructor(url, handlers) {
    this.url = url;
    this.handlers = handlers;
    this.ws = null;
    this.reconnectDelay = 1000;
    this.maxReconnectDelay = 10000;
    this.connected = false;
    this.connect();
  }

  connect() {
    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      this.connected = true;
      this.reconnectDelay = 1000;
      if (this.handlers.onConnect) this.handlers.onConnect();
    };

    this.ws.onclose = () => {
      this.connected = false;
      if (this.handlers.onDisconnect) this.handlers.onDisconnect();
      setTimeout(() => this.connect(), this.reconnectDelay);
      this.reconnectDelay = Math.min(this.reconnectDelay * 1.5, this.maxReconnectDelay);
    };

    this.ws.onerror = () => {
      this.ws.close();
    };

    this.ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        const handler = this.handlers[msg.type];
        if (handler) {
          handler(msg.payload || null);
        }
      } catch (e) {
        console.error('WS message parse error:', e);
      }
    };
  }

  send(type, payload) {
    if (!this.connected) return;
    const msg = { type };
    if (payload !== undefined) {
      msg.payload = payload;
    }
    this.ws.send(JSON.stringify(msg));
  }

  pause() { this.send('pause'); }
  resume() { this.send('resume'); }
  reset() { this.send('reset'); }
  saveBrain(name) { this.send('save_brain', { name }); }
  loadBrain(name) { this.send('load_brain', { name }); }
}
