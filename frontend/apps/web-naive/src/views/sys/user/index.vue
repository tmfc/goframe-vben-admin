<script lang="ts" setup>
import type { PaginationProps } from 'naive-ui';
import type { User } from '@vben/types';

import { h, onMounted, ref } from 'vue';

import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NSpace,
  NTag,
  useDialog,
} from 'naive-ui';

import { $t } from '#/locales';

import { getDeptTree } from '#/api/sys/dept';
import { deleteUser, getUserList } from '#/api/sys/user';
import UserFormModal from './modules/form.vue';

const dialog = useDialog();

function createColumns({
  edit,
  del,
}: {
  del: (id: string) => void;
  edit: (record: User) => void;
}) {
  return [
    {
      title: $t('system.user.columns.username'),
      key: 'username',
    },
    {
      title: $t('system.user.columns.realName'),
      key: 'realName',
    },
    {
      title: $t('system.user.columns.roles'),
      key: 'roles',
      render(row: User) {
        return formatRolesDisplay((row as any).roles);
      },
    },
    {
      title: $t('system.user.columns.deptId'),
      key: 'deptId',
      render(row: User) {
        const deptId = row?.deptId;
        if (!deptId) {
          return '';
        }
        return deptNameMap.value.get(String(deptId)) || String(deptId);
      },
    },
    {
      title: $t('system.user.columns.status'),
      key: 'status',
      render(row: User) {
        const isActive = row.status === 1;
        return h(
          NTag,
          { type: isActive ? 'success' : 'default', size: 'small' },
          {
            default: () =>
              isActive
                ? $t('system.user.status.enabled')
                : $t('system.user.status.disabled'),
          },
        );
      },
    },
    {
      title: $t('system.user.columns.actions'),
      key: 'actions',
      render(row: User) {
        return h('div', [
          h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              style: 'margin-right: 8px;',
              onClick: () => edit(row),
            },
            { default: () => $t('common.edit') },
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              onClick: () => del(row.id),
            },
            { default: () => $t('common.delete') },
          ),
        ]);
      },
    },
  ];
}

const data = ref<User[]>([]);
const deptNameMap = ref(new Map<string, string>());
const pagination = ref<PaginationProps>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
});
const loading = ref(false);

const showModal = ref(false);
const editingRecord = ref<User | null>(null);

const formValue = ref({
  username: '',
});

function formatRolesDisplay(value: unknown) {
  if (!value) return '';
  if (Array.isArray(value)) {
    return value.join(', ');
  }
  if (typeof value === 'string') {
    const trimmed = value.trim();
    if (!trimmed) return '';
    if (trimmed.startsWith('[')) {
      try {
        const parsed = JSON.parse(trimmed);
        if (Array.isArray(parsed)) {
          return parsed.join(', ');
        }
      } catch {
        return trimmed;
      }
    }
    return trimmed;
  }
  return String(value);
}

function handleCreate() {
  editingRecord.value = null;
  showModal.value = true;
}

function handleEdit(record: User) {
  editingRecord.value = record;
  showModal.value = true;
}

function handleDelete(id: string) {
  dialog.warning({
    title: $t('system.user.dialog.deleteTitle'),
    content: $t('system.user.dialog.deleteConfirm'),
    positiveText: $t('common.confirm'),
    negativeText: $t('common.cancel'),
    onPositiveClick: async () => {
      try {
        await deleteUser(id);
        fetchData();
      } catch (error) {
        console.error(error);
      }
    },
  });
}

function handleSearch() {
  pagination.value.page = 1;
  fetchData();
}

function handlePageChange(page: number) {
  pagination.value.page = page;
  fetchData();
}

async function fetchData() {
  loading.value = true;
  try {
    const response = await getUserList({
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
      username: formValue.value.username,
    });
    data.value = response.items;
    pagination.value.itemCount = response.total;
  } catch (error) {
    console.error(error);
  } finally {
    loading.value = false;
  }
}

function flattenDeptTree(
  list: Array<{ id: string; name: string; children?: any[] }> = [],
) {
  const map = new Map<string, string>();
  const walk = (items: any[]) => {
    for (const item of items || []) {
      if (item?.id) {
        map.set(String(item.id), item.name ?? '');
      }
      if (item?.children?.length) {
        walk(item.children);
      }
    }
  };
  walk(list);
  return map;
}

const columns = createColumns({
  edit: handleEdit,
  del: handleDelete,
});

onMounted(() => {
  getDeptTree()
    .then((res) => {
      deptNameMap.value = flattenDeptTree(res?.list || []);
    })
    .catch((error) => {
      console.error(error);
    });
  fetchData();
});
</script>

<template>
  <div class="user-page">
    <NCard :title="$t('system.user.title')" size="small">
      <div class="user-toolbar">
        <NForm
          inline
          :model="formValue"
          label-placement="left"
          label-width="auto"
          @submit.prevent="handleSearch"
        >
          <NFormItem :label="$t('system.user.filters.username')">
            <NInput
              v-model:value="formValue.username"
              :placeholder="$t('system.user.filters.searchByUsername')"
            />
          </NFormItem>
          <NFormItem>
            <NSpace>
              <NButton type="primary" attr-type="submit">
                {{ $t('common.search') }}
              </NButton>
            </NSpace>
          </NFormItem>
        </NForm>
        <NButton type="primary" @click="handleCreate">
          {{ $t('system.user.actions.create') }}
        </NButton>
      </div>
      <NDataTable
        :columns="columns"
        :data="data"
        :pagination="pagination"
        :loading="loading"
        @update:page="handlePageChange"
      />
    </NCard>
    <UserFormModal
      v-model:show="showModal"
      :record="editingRecord"
      @success="fetchData"
    />
  </div>
</template>

<style scoped>
.user-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.user-toolbar {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
}
</style>
