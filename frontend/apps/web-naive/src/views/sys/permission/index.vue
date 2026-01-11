<script lang="ts" setup>
import type { DataTableColumns } from 'naive-ui';

import { h, onMounted, reactive, ref } from 'vue';

import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NPopconfirm,
  NSelect,
  NSpace,
  NTag,
  useDialog,
} from 'naive-ui';

import { $t } from '#/locales';

import {
  deletePermission,
  getPermissionList,
  type PermissionItem,
} from '#/api/sys/permission';
import { collectExpandedKeys, listToTree } from '#/utils/tree';
import PermissionFormModal from './modules/form.vue';

const dialog = useDialog();
const loading = ref(false);
const showModal = ref(false);
const editingRecord = ref<PermissionItem | null>(null);

const filters = reactive({
  name: '',
  status: null as null | number,
});

const statusOptions = [
  { label: $t('system.permission.status.enabled'), value: 1 },
  { label: $t('system.permission.status.disabled'), value: 0 },
];

const data = ref<any[]>([]);
const rawList = ref<any[]>([]);
const defaultExpandedRowKeys = ref<string[]>([]);

const rowKey = (row: any) => row.id;

function buildTreeFromList(list: any[]) {
  const tree = listToTree(list);
  const expandedKeys = collectExpandedKeys(tree, 2);
  return { tree, expandedKeys };
}

const columns = reactive<DataTableColumns<any>>([
  {
    title: $t('system.permission.columns.name'),
    key: 'name',
    render(row) {
      return $t(row.name);
    },
  },
  {
    title: $t('system.permission.columns.description'),
    key: 'description',
    render(row) {
      return parsePermissionDescription(row.description);
    },
  },
  {
    title: $t('system.permission.columns.status'),
    key: 'status',
    render(row) {
      const isActive = row.status === 1;
      return h(
        NTag,
        { type: isActive ? 'success' : 'default', size: 'small' },
        {
          default: () =>
            isActive
              ? $t('system.permission.status.enabled')
              : $t('system.permission.status.disabled'),
        },
      );
    },
  },
  {
    title: $t('system.permission.columns.actions'),
    key: 'action',
    render(row) {
      return h(
        NSpace,
        { size: 8 },
        {
          default: () => [
            h(
              NButton,
              {
                size: 'small',
                type: 'primary',
                secondary: true,
                onClick: () => openEdit(row),
              },
              { default: () => $t('common.edit') },
            ),
            h(
              NPopconfirm,
              { onPositiveClick: () => handleDelete(row) },
              {
                trigger: () =>
                  h(
                    NButton,
                    { size: 'small', type: 'error', secondary: true },
                    { default: () => $t('common.delete') },
                  ),
                default: () => $t('system.permission.dialog.deleteConfirm'),
              },
            ),
          ],
        },
      );
    },
  },
]);

function parsePermissionDescription(desc: string): string {
  if (!desc) return '-';
  const match = desc.match(
    /^(菜单权限:|按钮权限:|Menu Permission:|Button Permission:)(.+)$/,
  );
  if (match) {
    const prefix = match[1];
    const i18nKey = match[2];
    if (i18nKey) {
      return `${prefix}${$t(i18nKey)}`;
    }
  }
  return desc;
}

async function fetchPermissionList() {
  try {
    loading.value = true;
    const res = await getPermissionList({
      name: filters.name || undefined,
      status: filters.status === null ? undefined : String(filters.status),
    });
    const list = res?.items || [];
    rawList.value = list;
    const { tree, expandedKeys } = buildTreeFromList(list);
    data.value = tree;
    defaultExpandedRowKeys.value = expandedKeys;
  } finally {
    loading.value = false;
  }
}

function handleSearch() {
  fetchPermissionList();
}

function resetFilters() {
  filters.name = '';
  filters.status = null;
  fetchPermissionList();
}

function openCreate() {
  editingRecord.value = null;
  showModal.value = true;
}

function openEdit(row: PermissionItem) {
  editingRecord.value = row;
  showModal.value = true;
}

function hasChildren(row: any) {
  const id = String(row.id);
  return rawList.value.some((item) => String(item.parentId || '') === id);
}

async function handleDelete(row: any) {
  if (hasChildren(row)) {
    dialog.warning({
      title: $t('system.permission.dialog.deleteTitle'),
      content: $t('system.permission.dialog.deleteHasChildren'),
      positiveText: $t('common.confirm'),
    });
    return;
  }
  try {
    await deletePermission(row.id);
    fetchPermissionList();
  } catch (error: any) {
    dialog.warning({
      title: $t('system.permission.dialog.deleteTitle'),
      content: error?.message || $t('system.permission.dialog.deleteFailed'),
      positiveText: $t('common.confirm'),
    });
  }
}

onMounted(() => {
  fetchPermissionList();
});
</script>

<template>
  <div class="permission-page">
    <NCard :title="$t('system.permission.title')" size="small">
      <NForm inline :model="filters" label-placement="left" label-width="auto">
        <NFormItem :label="$t('system.permission.filters.name')">
          <NInput
            v-model:value="filters.name"
            :placeholder="$t('system.permission.filters.searchByName')"
            clearable
          />
        </NFormItem>
        <NFormItem :label="$t('system.permission.filters.status')">
          <NSelect
            v-model:value="filters.status"
            :options="statusOptions"
            :placeholder="$t('system.permission.filters.all')"
            clearable
          />
        </NFormItem>
        <NFormItem>
          <NSpace>
            <NButton type="primary" @click="handleSearch">
              {{ $t('common.search') }}
            </NButton>
            <NButton @click="resetFilters">{{ $t('common.reset') }}</NButton>
          </NSpace>
        </NFormItem>
      </NForm>

      <div class="permission-toolbar">
        <NButton type="primary" @click="openCreate">
          {{ $t('system.permission.actions.create') }}
        </NButton>
      </div>

      <NDataTable
        :columns="columns"
        :data="data"
        :loading="loading"
        :bordered="false"
        :row-key="rowKey"
        :default-expanded-row-keys="defaultExpandedRowKeys"
      />
    </NCard>

    <PermissionFormModal
      v-model:show="showModal"
      :record="editingRecord"
      @success="fetchPermissionList"
    />
  </div>
</template>

<style scoped>
.permission-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.permission-toolbar {
  display: flex;
  justify-content: flex-end;
  margin: 12px 0;
}
</style>