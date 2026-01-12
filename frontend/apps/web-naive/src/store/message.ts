import { ref } from 'vue';

import { useAccessStore } from '@vben/stores';

import { defineStore } from 'pinia';

import { notification } from '#/adapter/naive';
import { WebSocketClient } from '#/utils/websocket';

export const useMessageStore = defineStore('message', () => {
  const accessStore = useAccessStore();
  const connected = ref(false);
  const messages = ref<any[]>([]);
  const unreadCount = ref(0);
  let client: WebSocketClient | null = null;

  function connect() {
    const token = accessStore.accessToken;
    if (!token) return;

    const apiUrl = import.meta.env.VITE_GLOB_API_URL || '/api';
    let wsUrl = '';

    if (apiUrl.startsWith('http')) {
      wsUrl = `${apiUrl.replace(/^http/, 'ws')}/ws`;
    } else {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = window.location.host;
      wsUrl = `${protocol}//${host}${apiUrl}/ws`;
    }

    // Append token
    wsUrl = `${wsUrl}?token=${token}`;

    if (client) {
      client.disconnect();
    }

    client = new WebSocketClient(wsUrl, {
      onClose: () => {
        connected.value = false;
      },
      onMessage: (event) => {
        try {
          const data = JSON.parse(event.data);
          messages.value.unshift(data);
          unreadCount.value++;

          notification.info({
            content: 'New Message',
            description: data.title || 'You have a new message',
            duration: 3000,
          });
        } catch (error) {
          console.error('Failed to parse message', error);
        }
      },
      onOpen: () => {
        connected.value = true;
      },
    });

    client.connect();
  }

  function disconnect() {
    client?.disconnect();
    client = null;
    connected.value = false;
  }

  return {
    connect,
    connected,
    disconnect,
    messages,
    unreadCount,
  };
});
