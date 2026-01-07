<script lang="ts" setup>
import type { PaginationProps } from 'naive-ui';

import type { Recordable } from '@vben/types';

import type { Dept } from '#/api/sys/dept';

import { h, onMounted, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '#/locales';

import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NSpace,
  useDialog,
} from 'naive-ui';

import { deleteDept, getDeptList } from '#/api/sys/dept';

import DeptFormModal from './modules/form.vue';

const dialog = useDialog();

function createColumns({
  edit,
  del,
}: {
  del: (id: string) => void;
  edit: (record: Dept) => void;
}) {
  return [
    {
      title: $t('system.dept.columns.name'),
      key: 'name',
    },
    {
      title: $t('system.dept.columns.status'),
      key: 'status',
      render(row: Dept) {
        return row.status === 1
          ? $t('system.dept.status.enabled')
          : $t('system.dept.status.disabled');
      },
    },
    {
      title: $t('system.dept.columns.order'),
      key: 'order',
    },
    {
      title: $t('system.dept.columns.actions'),
      key: 'actions',
      render(row: Dept) {
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

const data = ref<Dept[]>([]);
const pagination = ref<PaginationProps>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
});
const loading = ref(false);

const [Drawer, drawerApi] = useVbenDrawer({
  connectedComponent: DeptFormModal,
});

const formValue = ref({
  name: '',
});

function handleCreate() {
  drawerApi.open();
}

function handleEdit(record: Recordable) {
  drawerApi.open(record);
}

function handleDelete(id: string) {
  dialog.warning({
    title: $t('system.dept.dialog.deleteTitle'),
    content: $t('system.dept.dialog.deleteConfirm'),
    positiveText: $t('common.confirm'),
    negativeText: $t('common.cancel'),
    onPositiveClick: async () => {
      try {
        await deleteDept(id);
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
    const response = await getDeptList({
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
      name: formValue.value.name,
    });
    data.value = response.list;
    pagination.value.itemCount = response.total;
  } catch (error) {
    console.error(error);
  } finally {
    loading.value = false;
  }
}

const columns = createColumns({
  edit: handleEdit,
  del: handleDelete,
});

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="dept-page">
    <NCard :title="$t('system.dept.title')" size="small">
      <div class="dept-toolbar">
        <NForm inline :model="formValue" @submit.prevent="handleSearch">
          <NFormItem :label="$t('system.dept.filters.name')">
            <NInput
              v-model:value="formValue.name"
              :placeholder="$t('system.dept.filters.searchByName')"
            />
          </NFormItem>
          <NFormItem>
            <NSpace>
              <NButton type="primary" attr-type="submit">
                {{ $t('common.search') }}
              </NButton>
              <NButton type="primary" @click="handleCreate">
                {{ $t('system.dept.actions.create') }}
              </NButton>
            </NSpace>
          </NFormItem>
        </NForm>
      </div>
      <NDataTable
        :columns="columns"
        :data="data"
        :pagination="pagination"
        :loading="loading"
        @update:page="handlePageChange"
      />
    </NCard>
    <Drawer @success="fetchData" />
  </div>
</template>

<style scoped>
.dept-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.dept-toolbar {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
}
</style>
