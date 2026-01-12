import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { WebSocketClient } from '../websocket';

describe('WebSocketClient', () => {
  let client: WebSocketClient;
  const url = 'ws://localhost:1234';

  beforeEach(() => {
    vi.useFakeTimers();
    (global as any).WebSocket = vi.fn(() => ({
      close: vi.fn(),
      send: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      readyState: 1, 
    }));
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.restoreAllMocks();
  });

  it('should connect immediately', () => {
    client = new WebSocketClient(url);
    client.connect();
    expect(global.WebSocket).toHaveBeenCalledWith(url);
  });

  it('should not reconnect if manually closed', () => {
    client = new WebSocketClient(url);
    client.connect();
    client.disconnect();
    
    const wsInstance = (global.WebSocket as any).mock.results[0].value;
    
    // Trigger close event as if the server ack'd the close
    if (wsInstance.onclose) {
        wsInstance.onclose({ code: 1000 } as CloseEvent);
    }

    vi.advanceTimersByTime(5000);
    expect(global.WebSocket).toHaveBeenCalledTimes(1);
  });

  it('should reconnect on unexpected close', () => {
    client = new WebSocketClient(url, { reconnectInterval: 100, maxRetries: 2 });
    client.connect();
    
    const wsInstance = (global.WebSocket as any).mock.results[0].value;
    
    // Simulate unexpected close
    if (wsInstance.onclose) {
        wsInstance.onclose({ code: 1006 } as CloseEvent);
    }
    
    // Advance time to trigger reconnect
    vi.advanceTimersByTime(100);
    
    expect(global.WebSocket).toHaveBeenCalledTimes(2);
  });
});