<script lang="ts" setup>
import type { PaginationProps } from 'naive-ui';
import type { Role } from '#/api/sys/role';

import { h, onMounted, ref } from 'vue';

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

import { getRoleList, deleteRole } from '#/api/sys/role';
import RoleFormModal from './modules/form.vue';

const dialog = useDialog();

function createColumns({
  edit,
  del,
}: {
  del: (id: string) => void;
  edit: (record: Role) => void;
}) {
  return [
    {
      title: 'Role Name',
      key: 'name',
    },
    {
      title: 'Description',
      key: 'description',
    },
    {
      title: 'Status',
      key: 'status',
      render(row: Role) {
        return row.status === 1 ? 'Enabled' : 'Disabled';
      },
    },
    {
      title: 'Actions',
      key: 'actions',
      render(row: Role) {
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

const data = ref<Role[]>([]);
const pagination = ref<PaginationProps>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
});
const loading = ref(false);

const showModal = ref(false);
const editingRecord = ref<Role | null>(null);

const formValue = ref({
  name: '',
});

function handleCreate() {
  editingRecord.value = null;
  showModal.value = true;
}

function handleEdit(record: Role) {
  editingRecord.value = record;
  showModal.value = true;
}

function handleDelete(id: string) {
  dialog.warning({
    title: 'Delete Role',
    content: 'Are you sure you want to delete this role?',
    positiveText: 'Confirm',
    negativeText: 'Cancel',
    onPositiveClick: async () => {
      try {
        await deleteRole(id);
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
    const response = await getRoleList({
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
      name: formValue.value.name,
    });
    data.value = response.items;
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
  <div class="role-page">
    <NCard title="Role Management" size="small">
      <div class="role-toolbar">
        <NForm inline :model="formValue" @submit.prevent="handleSearch">
          <NFormItem label="Role Name">
            <NInput
              v-model:value="formValue.name"
              placeholder="Search by role name"
            />
          </NFormItem>
          <NFormItem>
            <NSpace>
              <NButton type="primary" attr-type="submit">
                Search
              </NButton>
              <NButton type="primary" @click="handleCreate">
                Create
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
    <RoleFormModal
      v-model:show="showModal"
      :record="editingRecord"
      @success="fetchData"
    />
  </div>
</template>

<style scoped>
.role-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.role-toolbar {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
}
</style>