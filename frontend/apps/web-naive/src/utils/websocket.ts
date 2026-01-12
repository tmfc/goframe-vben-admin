type OpenHandler = (event: Event) => void;
type MessageHandler = (event: MessageEvent) => void;
type ErrorHandler = (event: Event) => void;
type CloseHandler = (event: CloseEvent) => void;

export interface WebSocketClientOptions {
  reconnectInterval?: number;
  maxRetries?: number;
  onOpen?: OpenHandler;
  onMessage?: MessageHandler;
  onError?: ErrorHandler;
  onClose?: CloseHandler;
}

export class WebSocketClient {
  public url: string;
  private socket: WebSocket | null = null;
  private options: WebSocketClientOptions;
  private retryCount = 0;
  private isManuallyClosed = false;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;

  constructor(url: string, options: WebSocketClientOptions = {}) {
    this.url = url;
    this.options = {
      reconnectInterval: 5000,
      maxRetries: -1,
      ...options,
    };
  }

  connect() {
    this.isManuallyClosed = false;
    this.initSocket();
  }

  private initSocket() {
    this.socket = new WebSocket(this.url);

    this.socket.onopen = (event) => {
      this.retryCount = 0;
      this.options.onOpen?.(event);
    };

    this.socket.onmessage = (event) => {
      this.options.onMessage?.(event);
    };

    this.socket.onclose = (event) => {
      this.options.onClose?.(event);
      if (!this.isManuallyClosed) {
        this.reconnect();
      }
    };

    this.socket.onerror = (event) => {
      this.options.onError?.(event);
    };
  }

  private reconnect() {
    if (this.options.maxRetries !== -1 && this.retryCount >= (this.options.maxRetries || 0)) {
      return;
    }

    if (this.reconnectTimer) return;

    this.reconnectTimer = setTimeout(() => {
      this.retryCount++;
      this.reconnectTimer = null;
      this.initSocket();
    }, this.options.reconnectInterval);
  }

  disconnect() {
    this.isManuallyClosed = true;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }

  send(data: any) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(data);
    }
  }
}
