<template>
  <div>
    <n-data-table
      :columns="columns"
      :data="data"
      :pagination="pagination"
      :bordered="false"
    />
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, reactive, ref } from 'vue';
  import { NDataTable } from 'naive-ui';
  import { getMenuList } from '~/api/sys/menu';

  const loading = ref(false);

  const columns = reactive([
    {
      title: 'Name',
      key: 'name',
    },
    {
      title: 'Path',
      key: 'path',
    },
    {
      title: 'Component',
      key: 'component',
    },
    {
      title: 'Type',
      key: 'type',
    },
    {
      title: 'Status',
      key: 'status',
    },
    {
      title: 'Action',
      key: 'action',
      render() {
        return 'TODO';
      },
    },
  ]);

  const data = ref([]);

  const pagination = reactive({
    page: 1,
    pageSize: 10,
    itemCount: 0,
    onChange: (page) => {
      pagination.page = page;
      fetchMenuList();
    },
    onUpdatePageSize: (pageSize) => {
      pagination.pageSize = pageSize;
      pagination.page = 1;
      fetchMenuList();
    },
  });

  async function fetchMenuList() {
    try {
      loading.value = true;
      const res = await getMenuList({
        page: pagination.page,
        pageSize: pagination.pageSize,
      });
      data.value = res.list;
      pagination.itemCount = res.total;
    } finally {
      loading.value = false;
    }
  }

  onMounted(() => {
    fetchMenuList();
  });
</script>
