<template>
  <div>
    <n-data-table
      :columns="columns"
      :data="data"
      :pagination="pagination"
      :loading="loading"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { NDataTable } from 'naive-ui';
import { getDeptList } from '#/api/sys/dept';

const columns = [
  {
    title: 'Name',
    key: 'name',
  },
  {
    title: 'Status',
    key: 'status',
  },
  {
    title: 'Order',
    key: 'order',
  },
];

const data = ref([]);
const pagination = ref({
  page: 1,
  pageSize: 10,
  itemCount: 0,
});
const loading = ref(false);

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