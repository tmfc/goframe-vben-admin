<template>
  <div>
    <div class="mb-2">
      <n-button type="primary" @click="handleCreate">
        Create
      </n-button>
    </div>
    <n-data-table
      :columns="columns"
      :data="data"
      :pagination="pagination"
      :loading="loading"
    />
    <Drawer @success="fetchData" />
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, h } from 'vue';
import { NDataTable, NButton, useDialog } from 'naive-ui';
import { getDeptList, deleteDept } from '#/api/sys/dept';
import Drawer from './modules/form.vue';
import { useDrawer } from '@vben/common-ui';

const dialog = useDialog();

const columns = [
  {
    title: 'Name',
    key: 'name',
  },
  {
    title: 'Status',
    key: 'status',
    render(row) {
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
    render(row) {
      return h('div', [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-right: 8px;',
            onClick: () => handleEdit(row),
          },
          { default: () => 'Edit' }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            onClick: () => handleDelete(row.id),
          },
          { default: () => 'Delete' }
        ),
      ]);
    },
  },
];

const data = ref([]);
const pagination = ref({
  page: 1,
  pageSize: 10,
  itemCount: 0,
});
const loading = ref(false);

const [DrawerRef, drawerApi] = useDrawer();

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

async function fetchData() {
  loading.value = true;
  try {
    const response = await getDeptList({
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
    });
    data.value = response.list;
    pagination.value.itemCount = response.total;
  } catch (error) {
    console.error(error);
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchData();
});
</script>
