import { computed } from 'vue';
import { useMessageStore } from '#/store/message';

export function useNotification() {
  const messageStore = useMessageStore();

  const notifications = computed(() => {
    return messageStore.messages.map((msg) => ({
      id: msg.id,
      title: msg.title,
      message: msg.content || msg.message,
      date: msg.created_at,
      isRead: !!msg.read_at || !!msg.read,
      avatar: msg.sender_avatar || 'https://avatar.vercel.sh/vercel.svg?text=VB',
    }));
  });

  const showDot = computed(() => messageStore.unreadCount > 0);

  function handleNoticeClear() {
    messageStore.messages = [];
    messageStore.unreadCount = 0;
  }

  function markRead(id: number | string) {
    const msg = messageStore.messages.find((m) => m.id === id);
    if (msg && !msg.read && !msg.read_at) {
      msg.read = true;
      if (messageStore.unreadCount > 0) {
        messageStore.unreadCount--;
      }
    }
  }

  function handleMakeAll() {
    messageStore.messages.forEach((m) => {
      m.read = true;
    });
    messageStore.unreadCount = 0;
  }

  function remove(id: number | string) {
    const index = messageStore.messages.findIndex((m) => m.id === id);
    if (index !== -1) {
      const msg = messageStore.messages[index];
      if (!msg.read && !msg.read_at) {
        if (messageStore.unreadCount > 0) {
          messageStore.unreadCount--;
        }
      }
      messageStore.messages.splice(index, 1);
    }
  }

  return {
    handleMakeAll,
    handleNoticeClear,
    markRead,
    notifications,
    remove,
    showDot,
  };
}
