import { describe, it, expect, vi, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { reactive } from 'vue';
import { useNotification } from '../useNotification';
import { useMessageStore } from '#/store/message';

// Mock dependencies
vi.mock('#/store/message', () => ({
  useMessageStore: vi.fn(),
}));

describe('useNotification', () => {
  let storeMock: any;

  beforeEach(() => {
    setActivePinia(createPinia());
    storeMock = reactive({
      messages: [],
      unreadCount: 0,
      connect: vi.fn(),
    });
    (useMessageStore as any).mockReturnValue(storeMock);
  });

  it('should map store messages to notifications', () => {
    storeMock.messages = [
      { id: 1, title: 'Title', content: 'Content', created_at: '2023-01-01T00:00:00Z', read_at: null }
    ];

    const { notifications } = useNotification();

    expect(notifications.value).toHaveLength(1);
    expect(notifications.value[0]).toMatchObject({
      id: 1,
      title: 'Title',
      message: 'Content',
      isRead: false,
    });
  });

  it('should handle read status correctly', () => {
    storeMock.messages = [
      { id: 1, title: 'Title', content: 'Content', created_at: '2023-01-01T00:00:00Z', read_at: '2023-01-02' }
    ];

    const { notifications } = useNotification();
    expect(notifications.value[0].isRead).toBe(true);
  });
});
