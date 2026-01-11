<script lang="ts" setup>
import type { PaginationProps } from 'naive-ui';
import type { TenantModel } from '#/api/sys/tenant';

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

import { deleteTenant, getTenantList } from '#/api/sys/tenant';
import TenantFormModal from './modules/form.vue';

const dialog = useDialog();

function createColumns({
  edit,
  del,
}: {
  del: (id: string | number) => void;
  edit: (record: TenantModel) => void;
}) {
  return [
    {
      title: $t('system.tenant.columns.name'),
      key: 'name',
    },
    {
      title: $t('system.tenant.columns.code'),
      key: 'code',
    },
    {
      title: $t('system.tenant.columns.contactName'),
      key: 'contactName',
    },
    {
      title: $t('system.tenant.columns.contactMobile'),
      key: 'contactMobile',
    },
    {
      title: $t('system.tenant.columns.status'),
      key: 'status',
      render(row: TenantModel) {
        const isActive = row.status === 1;
        return h(
          NTag,
          { type: isActive ? 'success' : 'default', size: 'small' },
          {
            default: () =>
              isActive
                ? $t('system.tenant.status.enabled')
                : $t('system.tenant.status.disabled'),
          },
        );
      },
    },
    {
      title: $t('system.tenant.columns.expireAt'),
      key: 'expireAt',
    },
    {
      title: $t('system.tenant.columns.actions'),
      key: 'actions',
      render(row: TenantModel) {
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
              onClick: () => del(row.id!),
            },
            { default: () => $t('common.delete') },
          ),
        ]);
      },
    },
  ];
}

const data = ref<TenantModel[]>([]);
const pagination = ref<PaginationProps>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
});
const loading = ref(false);

const showModal = ref(false);
const editingRecord = ref<TenantModel | null>(null);

const formValue = ref({
  name: '',
  code: '',
});

function handleCreate() {
  editingRecord.value = null;
  showModal.value = true;
}

function handleEdit(record: TenantModel) {
  editingRecord.value = record;
  showModal.value = true;
}

function handleDelete(id: string | number) {
  dialog.warning({
    title: $t('system.tenant.dialog.deleteTitle'),
    content: $t('system.tenant.dialog.deleteConfirm'),
    positiveText: $t('common.confirm'),
    negativeText: $t('common.cancel'),
    onPositiveClick: async () => {
      try {
        await deleteTenant(id);
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
    const response = await getTenantList({
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
      name: formValue.value.name,
      code: formValue.value.code,
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
  <div class="tenant-page">
    <NCard :title="$t('system.tenant.title')" size="small">
      <div class="tenant-toolbar">
        <NForm
          inline
          :model="formValue"
          label-placement="left"
          label-width="auto"
          @submit.prevent="handleSearch"
        >
          <NFormItem :label="$t('system.tenant.filters.name')">
            <NInput
              v-model:value="formValue.name"
              :placeholder="$t('system.tenant.filters.searchByName')"
            />
          </NFormItem>
          <NFormItem :label="$t('system.tenant.filters.code')">
            <NInput
              v-model:value="formValue.code"
              placeholder="请输入租户代码"
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
          {{ $t('system.tenant.actions.create') }}
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
    <TenantFormModal
      v-model:show="showModal"
      :record="editingRecord"
      @success="fetchData"
    />
  </div>
</template>

<style scoped>
.tenant-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.tenant-toolbar {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
}
</style>
