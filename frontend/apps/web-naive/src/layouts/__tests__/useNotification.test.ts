import { describe, it, expect, vi, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useNotification } from '../useNotification';
import { useMessageStore } from '#/store/message';
import * as messageApi from '#/api/sys/message';

// Mock dependencies
vi.mock('#/store/message', () => ({
  useMessageStore: vi.fn(),
}));

vi.mock('vue-router', () => ({
  useRouter: vi.fn(),
}));

vi.mock('#/api/sys/message', () => ({
  setMessageReadApi: vi.fn().mockResolvedValue({}),
}));

describe('useNotification', () => {
  let storeMock: any;
  let routerMock: any;

  beforeEach(() => {
    setActivePinia(createPinia());
    storeMock = reactive({
      messages: [],
      unreadCount: 0,
      connect: vi.fn(),
    });
    (useMessageStore as any).mockReturnValue(storeMock);

    routerMock = {
      push: vi.fn(),
    };
    (useRouter as any).mockReturnValue(routerMock);
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

  it('should navigate when jump_path is present', async () => {
    storeMock.messages = [
      { id: 1, title: 'Title', jump_path: '/test', jump_params: '{"a":1}', read_at: null }
    ];

    const { markRead } = useNotification();
    await markRead(1);

    expect(routerMock.push).toHaveBeenCalledWith({ path: '/test', query: { a: 1 } });
  });
});
