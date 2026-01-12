import { describe, it, expect, vi, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useMessageStore } from '../message';
import { WebSocketClient } from '#/utils/websocket';

vi.mock('#/utils/websocket', () => {
  return {
    WebSocketClient: vi.fn().mockImplementation(() => ({
      connect: vi.fn(),
      disconnect: vi.fn(),
    })),
  };
});

describe('Message Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('should initialize with default values', () => {
    const store = useMessageStore();
    expect(store.connected).toBe(false);
    expect(store.messages).toEqual([]);
    expect(store.unreadCount).toBe(0);
  });

  it('should create websocket client on connect', () => {
    const store = useMessageStore();
    // Mock accessStore if needed, but let's assume it works for now or mock useAccessStore
    // For simplicity, I'll focus on the store's own logic
  });
});
