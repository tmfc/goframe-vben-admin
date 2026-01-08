<script lang="ts" setup>
import type { Dept } from '#/api/sys/dept';

import { h, onMounted, ref } from 'vue';

import { $t } from '#/locales';

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

import { deleteDept, getDeptTree } from '#/api/sys/dept';

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
      title: $t('system.dept.columns.parent'),
      key: 'parentId',
      render(row: Dept) {
        const parentId = row?.parentId;
        if (!parentId) {
          return '';
        }
        return deptNameMap.value.get(String(parentId)) || String(parentId);
      },
    },
    {
      title: $t('system.dept.columns.status'),
      key: 'status',
      render(row: Dept) {
        const isActive = row.status === 1;
        return h(
          NTag,
          { type: isActive ? 'success' : 'default', size: 'small' },
          {
            default: () =>
              isActive
                ? $t('system.dept.status.enabled')
                : $t('system.dept.status.disabled'),
          },
        );
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
const deptNameMap = ref(new Map<string, string>());
const loading = ref(false);
const defaultExpandedRowKeys = ref<string[]>([]);

const showModal = ref(false);
const editingRecord = ref<Dept | null>(null);

const formValue = ref({
  name: '',
});

function handleCreate() {
  editingRecord.value = null;
  showModal.value = true;
}

function handleEdit(record: Dept) {
  editingRecord.value = record;
  showModal.value = true;
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
  fetchData();
}

const rowKey = (row: Dept) => row.id;

async function fetchData() {
  loading.value = true;
  try {
    const response = await getDeptTree({
      name: formValue.value.name,
    });
    const list = response.list || [];
    data.value = list;
    deptNameMap.value = flattenDeptTree(list);
    defaultExpandedRowKeys.value = collectExpandedKeys(list);
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

function collectExpandedKeys(
  list: Array<{ id: string; children?: any[] }> = [],
) {
  const keys: string[] = [];
  const walk = (items: Array<{ id: string; children?: any[] }>) => {
    for (const item of items || []) {
      if (item?.children?.length) {
        keys.push(String(item.id));
        walk(item.children);
      }
    }
  };
  walk(list);
  return keys;
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
        <NForm
          inline
          :model="formValue"
          label-placement="left"
          label-width="auto"
          @submit.prevent="handleSearch"
        >
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
            </NSpace>
          </NFormItem>
        </NForm>
        <NButton type="primary" @click="handleCreate">
          {{ $t('system.dept.actions.create') }}
        </NButton>
      </div>
      <NDataTable
        :columns="columns"
        :data="data"
        :loading="loading"
        :row-key="rowKey"
        :default-expanded-row-keys="defaultExpandedRowKeys"
      />
    </NCard>
    <DeptFormModal
      v-model:show="showModal"
      :record="editingRecord"
      @success="fetchData"
    />
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
