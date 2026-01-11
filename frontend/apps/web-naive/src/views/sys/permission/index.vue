<script lang="ts" setup>
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui';

import { computed, h, onMounted, reactive, ref } from 'vue';

import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NPopconfirm,
  NSelect,
  NSpace,
  NTag,
  NTreeSelect,
  useDialog,
} from 'naive-ui';

import { $t } from '#/locales';

import {
  createPermission,
  deletePermission,
  getPermissionList,
  updatePermission,
} from '#/api/sys/permission';

const dialog = useDialog();
const loading = ref(false);
const saving = ref(false);
const showModal = ref(false);
const editingId = ref<null | string>(null);
const formRef = ref<FormInst | null>(null);

const filters = reactive({
  name: '',
  status: null as null | number,
});

const form = reactive({
  name: '',
  description: '',
  parentId: null as null | string,
  status: 1,
});

const rules: FormRules = {
  name: [
    {
      required: true,
      message: $t('system.permission.validation.nameRequired'),
      trigger: 'blur',
    },
  ],
};

const statusOptions = [
  { label: $t('system.permission.status.enabled'), value: 1 },
  { label: $t('system.permission.status.disabled'), value: 0 },
];

const data = ref<any[]>([]);
const rawList = ref<any[]>([]);
const defaultExpandedRowKeys = ref<string[]>([]);
const treeOptions = ref<any[]>([]);

const rowKey = (row: any) => row.id;

const modalTitle = computed(() =>
  editingId.value
    ? $t('system.permission.actions.edit')
    : $t('system.permission.actions.create'),
);

function buildTreeFromList(list: any[]) {
  const nodeMap = new Map<string, any>();
  const roots: any[] = [];

  for (const item of list || []) {
    const id = String(item.id);
    nodeMap.set(id, { ...item, id, children: [] });
  }

  for (const item of list || []) {
    const id = String(item.id);
    const parentId = item.parentId ? String(item.parentId) : '';
    const node = nodeMap.get(id);
    if (!node) continue;
    if (!parentId || parentId === '0' || !nodeMap.has(parentId)) {
      roots.push(node);
    } else {
      nodeMap.get(parentId).children.push(node);
    }
  }

  const expandedKeys: string[] = [];
  const collectExpand = (nodes: any[], depth: number) => {
    for (const node of nodes) {
      if (depth <= 2) {
        expandedKeys.push(node.id);
      }
      if (node.children?.length) {
        collectExpand(node.children, depth + 1);
      }
    }
  };
  collectExpand(roots, 1);

  return { tree: roots, expandedKeys };
}

function buildTreeOptions(list: any[]) {
  const nodeMap = new Map<string, any>();
  const roots: any[] = [];

  for (const item of list || []) {
    const id = String(item.id);
    nodeMap.set(id, { key: id, label: item.name ?? id, children: [] });
  }

  for (const item of list || []) {
    const id = String(item.id);
    const parentId = item.parentId ? String(item.parentId) : '';
    const node = nodeMap.get(id);
    if (!node) continue;
    if (!parentId || parentId === '0' || !nodeMap.has(parentId)) {
      roots.push(node);
    } else {
      nodeMap.get(parentId).children.push(node);
    }
  }

  return roots;
}

const columns = reactive<DataTableColumns<any>>([
  {
    title: $t('system.permission.columns.name'),
    key: 'name',
    render(row) {
      // name 是 i18n key,直接翻译
      return $t(row.name);
    },
  },
  {
    title: $t('system.permission.columns.description'),
    key: 'description',
    render(row) {
      // description 是格式化字符串,需要解析
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

// 解析权限描述中的 i18n key
function parsePermissionDescription(desc: string): string {
  if (!desc) return '-';

  // 格式: "菜单权限:system.user.title" 或 "按钮权限:common.create"
  const match = desc.match(/^(菜单权限:|按钮权限:|Menu Permission:|Button Permission:)(.+)$/);
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
      page: 1,
      pageSize: 1000,
      name: filters.name || undefined,
      status: filters.status === null ? undefined : String(filters.status),
    });
    const list = res?.items || [];
    rawList.value = list;
    const { tree, expandedKeys } = buildTreeFromList(list);
    data.value = tree;
    defaultExpandedRowKeys.value = expandedKeys;
    treeOptions.value = [
      { key: '0', label: $t('system.permission.form.root') },
      ...buildTreeOptions(list),
    ];
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

function resetForm() {
  form.name = '';
  form.description = '';
  form.parentId = null;
  form.status = 1;
  formRef.value?.restoreValidation();
}

function openCreate() {
  editingId.value = null;
  resetForm();
  showModal.value = true;
}

function openEdit(row: any) {
  editingId.value = row.id;
  form.name = row.name ?? '';
  form.description = row.description ?? '';
  form.parentId = row.parentId ? String(row.parentId) : null;
  form.status = row.status ?? 1;
  showModal.value = true;
}

function closeModal() {
  showModal.value = false;
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

async function submitForm() {
  if (!formRef.value) return;
  try {
    await formRef.value.validate();
  } catch {
    return;
  }
  saving.value = true;
  try {
    const parentIdValue = form.parentId ? Number(form.parentId) : 0;
    const payload = {
      name: form.name,
      description: form.description,
      parent_id: Number.isFinite(parentIdValue) ? parentIdValue : 0,
      status: form.status,
    };
    await (editingId.value
      ? updatePermission(editingId.value, payload)
      : createPermission(payload));
    showModal.value = false;
    fetchPermissionList();
  } finally {
    saving.value = false;
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

    <NModal
      v-model:show="showModal"
      preset="dialog"
      :title="modalTitle"
      :mask-closable="false"
    >
      <NForm ref="formRef" :model="form" :rules="rules" label-placement="top">
        <NFormItem :label="$t('system.permission.form.name')" path="name">
          <NInput
            v-model:value="form.name"
            :placeholder="$t('system.permission.form.namePlaceholder')"
          />
        </NFormItem>
        <NFormItem :label="$t('system.permission.form.description')">
          <NInput
            v-model:value="form.description"
            :placeholder="$t('system.permission.form.descriptionPlaceholder')"
          />
        </NFormItem>
        <NFormItem :label="$t('system.permission.form.parent')">
          <NTreeSelect
            v-model:value="form.parentId"
            :options="treeOptions"
            key-field="key"
            label-field="label"
            children-field="children"
            clearable
          />
        </NFormItem>
        <NFormItem :label="$t('system.permission.form.status')">
          <NSelect v-model:value="form.status" :options="statusOptions" />
        </NFormItem>
      </NForm>
      <template #action>
        <NSpace>
          <NButton @click="closeModal">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" :loading="saving" @click="submitForm">
            {{ $t('common.confirm') }}
          </NButton>
        </NSpace>
      </template>
    </NModal>
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
