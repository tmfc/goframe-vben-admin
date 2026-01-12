import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useMessageStore } from '#/store/message';
import { setMessageReadApi } from '#/api/sys/message';

export function useNotification() {
  const messageStore = useMessageStore();
  const router = useRouter();

  const notifications = computed(() => {
    return messageStore.messages.map((msg) => ({
      id: msg.id,
      title: msg.title,
      message: msg.content || msg.message,
      date: msg.created_at,
      isRead: !!msg.read_at || !!msg.read,
      avatar: msg.sender_avatar || 'https://avatar.vercel.sh/vercel.svg?text=VB',
      jump_path: msg.jump_path,
      jump_params: msg.jump_params,
    }));
  });

  const showDot = computed(() => messageStore.unreadCount > 0);

  function handleNoticeClear() {
    messageStore.messages = [];
    messageStore.unreadCount = 0;
  }

  async function markRead(id: number | string) {
    const msg = messageStore.messages.find((m) => m.id === id);
    if (msg) {
      if (!msg.read && !msg.read_at) {
        msg.read = true;
        if (messageStore.unreadCount > 0) {
          messageStore.unreadCount--;
        }
        // Call API in background
        setMessageReadApi([Number(id)]).catch(console.error);
      }

      // Handle jump
      if (msg.jump_path) {
        try {
          const query = msg.jump_params ? JSON.parse(msg.jump_params) : {};
          router.push({ path: msg.jump_path, query });
        } catch (e) {
          console.error('Failed to parse jump_params', e);
          router.push({ path: msg.jump_path });
        }
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
