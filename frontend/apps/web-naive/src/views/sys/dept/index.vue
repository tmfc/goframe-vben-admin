<script lang="ts" setup>
import type { PaginationProps } from 'naive-ui';

import type { Recordable } from '@vben/types';

import type { Dept } from '#/api/sys/dept';

import { h, onMounted, ref } from 'vue';

import { useDrawer } from '@vben/common-ui';

import {
  NButton,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
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
      title: 'Name',
      key: 'name',
    },
    {
      title: 'Status',
      key: 'status',
      render(row: Dept) {
        return row.status === 1 ? 'Enabled' : 'Disabled';
      },
    },
    {
      title: 'Order',
      key: 'order',
    },
    {
      title: 'Actions',
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
            { default: () => 'Edit' },
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              onClick: () => del(row.id),
            },
            { default: () => 'Delete' },
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

const [, drawerApi] = useDrawer();

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
    title: 'Warning',
    content: 'Are you sure you want to delete this department?',
    positiveText: 'Yes',
    negativeText: 'No',
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
  <div>
    <div class="mb-2">
      <NForm inline :model="formValue" @submit.prevent="handleSearch">
        <NFormItem>
          <NInput v-model:value="formValue.name" placeholder="Search by Name" />
        </NFormItem>
        <NFormItem>
          <NButton type="primary" attr-type="submit"> Search </NButton>
        </NFormItem>
        <NFormItem>
          <NButton type="primary" @click="handleCreate"> Create </NButton>
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
    <DeptFormModal @success="fetchData" />
  </div>
</template>
