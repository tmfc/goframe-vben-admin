<script lang="ts" setup>
import type { PaginationProps } from 'naive-ui';
import { h, onMounted, ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import {
  NButton,
  NCard,
  NDataTable,
  NSpace,
  NTag,
  NInput,
  NSelect,
  NForm,
  NFormItem,
} from 'naive-ui';
import { $t } from '#/locales';
import { getMessageListApi, setMessageReadApi, setAllMessageReadApi, type MessageApi } from '#/api/sys/message';

const router = useRouter();
const loading = ref(false);
const data = ref<MessageApi.MessageItem[]>([]);
const total = ref(0);
const pagination = reactive<PaginationProps>({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.page = page;
    fetchData();
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
    fetchData();
  },
});

const searchParams = reactive({
  title: '',
  status: null as number | null,
});

async function fetchData() {
  try {
    loading.value = true;
    const res = await getMessageListApi({
      page: pagination.page || 1,
      size: pagination.pageSize || 10,
    });
    data.value = res.list;
    total.value = res.total;
    pagination.itemCount = res.total;
  } finally {
    loading.value = false;
  }
}

async function handleMarkRead(row: MessageApi.MessageItem) {
  if (!row.read_at) {
    await setMessageReadApi([row.id]);
    fetchData();
  }

  if (row.jump_path) {
    try {
      const query = row.jump_params ? JSON.parse(row.jump_params) : {};
      router.push({ path: row.jump_path, query });
    } catch (e) {
      console.error('Failed to parse jump_params', e);
      router.push({ path: row.jump_path });
    }
  }
}

async function handleMarkAllRead() {
  await setAllMessageReadApi();
  fetchData();
}

const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: $t('system.message.columns.title'),
    key: 'title',
    render(row: MessageApi.MessageItem) {
      return h(
        'span',
        {
          style: {
            fontWeight: row.read_at ? 'normal' : 'bold',
            cursor: 'pointer',
          },
          onClick: () => handleMarkRead(row),
        },
        row.title,
      );
    },
  },
  {
    title: $t('system.message.columns.content'),
    key: 'content',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: $t('system.message.columns.status'),
    key: 'read_at',
    width: 100,
    render(row: MessageApi.MessageItem) {
      return h(
        NTag,
        {
          type: row.read_at ? 'default' : 'error',
          size: 'small',
        },
        { default: () => (row.read_at ? $t('system.message.status.read') : $t('system.message.status.unread')) },
      );
    },
  },
  {
    title: $t('system.message.columns.createdAt'),
    key: 'created_at',
    width: 180,
  },
  {
    title: $t('system.message.columns.actions'),
    key: 'actions',
    width: 120,
    render(row: MessageApi.MessageItem) {
      return h(
        NSpace,
        {},
        {
          default: () => [
            h(
              NButton,
              {
                size: 'small',
                type: 'primary',
                ghost: true,
                disabled: !!row.read_at,
                onClick: () => handleMarkRead(row),
              },
              { default: () => $t('system.message.actions.markRead') },
            ),
          ],
        },
      );
    },
  },
];

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="p-4">
    <NCard :title="$t('system.message.title')" :bordered="false" class="shadow-sm">
      <template #header-extra>
        <NSpace>
          <NButton type="primary" @click="handleMarkAllRead">
            {{ $t('system.message.actions.markAllRead') }}
          </NButton>
          <NButton @click="fetchData">刷新</NButton>
        </NSpace>
      </template>

      <NDataTable
        remote
        :loading="loading"
        :columns="columns"
        :data="data"
        :pagination="pagination"
        :bordered="false"
      />
    </NCard>
  </div>
</template>
