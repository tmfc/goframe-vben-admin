<script lang="ts" setup>
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui';

import { computed, h, onMounted, reactive, ref } from 'vue';

import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NFormItemGi,
  NGrid,
  NInput,
  NInputNumber,
  NModal,
  NPopconfirm,
  NSelect,
  NSpace,
  NTag,
  NTreeSelect,
} from 'naive-ui';

import {
  createMenu,
  deleteMenu,
  getMenuList,
  updateMenu,
} from '#/api/sys/menu';
import { $t } from '#/locales';

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
  path: '',
  component: '',
  icon: '',
  type: 'menu',
  parentId: null as null | string,
  metaTitle: '',
  status: 1,
  order: 0,
});

const rules: FormRules = {
  name: [
    {
      required: true,
      message: $t('system.menu.validation.nameRequired'),
      trigger: 'blur',
    },
  ],
  type: [
    {
      required: true,
      message: $t('system.menu.validation.typeRequired'),
      trigger: 'change',
    },
  ],
};

const statusOptions = [
  { label: $t('system.menu.status.enabled'), value: 1 },
  { label: $t('system.menu.status.disabled'), value: 0 },
];

const typeOptions = [
  { label: $t('system.menu.types.menu'), value: 'menu' },
  { label: $t('system.menu.types.catalog'), value: 'catalog' },
  { label: $t('system.menu.types.link'), value: 'link' },
  { label: $t('system.menu.types.embedded'), value: 'embedded' },
];

const data = ref<any[]>([]);
const defaultExpandedRowKeys = ref<string[]>([]);
const menuTreeOptions = ref<any[]>([]);
const rawMenuList = ref<any[]>([]);

function getMeta(row: any) {
  if (!row) return null;
  const rawMeta = row.meta;
  if (typeof rawMeta === 'string') {
    try {
      return JSON.parse(rawMeta);
    } catch {
      return null;
    }
  }
  if (rawMeta && typeof rawMeta === 'object') return rawMeta;
  return null;
}

function getMenuTitle(row: any) {
  if (!row) return '';
  const meta = getMeta(row);
  const titleKey = meta?.title;
  if (titleKey) {
    return $t(titleKey);
  }
  return row.name ?? '';
}

function getTitleKey(row: any) {
  const meta = getMeta(row);
  return meta?.title ?? '';
}

function buildMenuTreeFromList(list: any[]) {
  const nodeMap = new Map<string, any>();
  const roots: any[] = [];

  for (const item of list || []) {
    const id = String(item.id);
    nodeMap.set(id, {
      id,
      label: getMenuTitle(item),
      children: [],
    });
  }

  for (const item of list || []) {
    const id = String(item.id);
    const parentId = item.parentId ? String(item.parentId) : '';
    const node = nodeMap.get(id);
    if (!node) {
      continue;
    }
    if (!parentId || parentId === '0' || !nodeMap.has(parentId)) {
      roots.push(node);
      continue;
    }
    nodeMap.get(parentId).children.push(node);
  }

  return roots;
}

const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  onChange: (page: number) => {
    pagination.page = page;
    fetchMenuList();
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
    fetchMenuList();
  },
});

const columns = reactive<DataTableColumns<any>>([
  {
    title: $t('system.menu.columns.name'),
    key: 'name',
    render(row) {
      return getMenuTitle(row);
    },
  },
  { title: $t('system.menu.columns.path'), key: 'path' },
  { title: $t('system.menu.columns.component'), key: 'component' },
  {
    title: $t('system.menu.columns.type'),
    key: 'type',
    render(row) {
      const option = typeOptions.find((item) => item.value === row.type);
      return option?.label ?? row.type ?? '';
    },
  },
  { title: $t('system.menu.columns.order'), key: 'order' },
  {
    title: $t('system.menu.columns.status'),
    key: 'status',
    render(row) {
      const isActive = row.status === 1;
      return h(
        NTag,
        { type: isActive ? 'success' : 'default', size: 'small' },
        {
          default: () =>
            isActive
              ? $t('system.menu.status.enabled')
              : $t('system.menu.status.disabled'),
        },
      );
    },
  },
  {
    title: $t('system.menu.columns.action'),
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
            row.type === 'menu'
              ? h(
                  NButton,
                  {
                    size: 'small',
                    tertiary: true,
                    onClick: () => openButtonModal(row),
                  },
                  { default: () => $t('system.menu.actions.manageButtons') },
                )
              : null,
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
                default: () => $t('system.menu.actions.deleteConfirm'),
              },
            ),
          ],
        },
      );
    },
  },
]);

const buttonColumns = reactive<DataTableColumns<any>>([
  { title: $t('system.menu.columns.name'), key: 'name' },
  {
    title: $t('system.menu.columns.titleKey'),
    key: 'titleKey',
    render(row) {
      const key = getTitleKey(row);
      if (!key) return '-';
      return `${$t(key)} (${key})`;
    },
  },
  {
    title: $t('system.menu.columns.path'),
    key: 'permission',
    render(row) {
      return row.path;
    },
  },
  { title: $t('system.menu.columns.order'), key: 'order' },
  {
    title: $t('system.menu.columns.status'),
    key: 'status',
    render(row) {
      const isActive = row.status === 1;
      return h(
        NTag,
        { type: isActive ? 'success' : 'default', size: 'small' },
        {
          default: () =>
            isActive ? $t('system.menu.status.enabled') : $t('system.menu.status.disabled'),
        },
      );
    },
  },
  {
    title: $t('system.menu.columns.action'),
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
                onClick: () => openEditButton(row),
              },
              { default: () => $t('common.edit') },
            ),
            h(
              NPopconfirm,
              { onPositiveClick: () => handleDeleteButton(row) },
              {
                trigger: () =>
                  h(
                    NButton,
                    { size: 'small', type: 'error', secondary: true },
                    { default: () => $t('common.delete') },
                  ),
                default: () => $t('system.menu.actions.deleteConfirm'),
              },
            ),
          ],
        },
      );
    },
  },
]);

const rowKey = (row: any) => row.id;

const modalTitle = computed(() =>
  editingId.value
    ? $t('system.menu.actions.edit')
    : $t('system.menu.actions.create'),
);

const buttonModal = reactive({
  show: false,
  menuId: '' as string,
  menuName: '',
});
const buttonList = ref<any[]>([]);
const buttonFormRef = ref<FormInst | null>(null);
const buttonSaving = ref(false);
const buttonEditingId = ref<null | string>(null);
const buttonForm = reactive({
  name: '',
  permission: '',
  metaTitle: '',
  status: 1,
  order: 0,
});

async function fetchMenuList() {
  try {
    loading.value = true;
    const res = await getMenuList({
      page: 1,
      pageSize: 1000,
      name: filters.name || undefined,
      status: filters.status === null ? undefined : String(filters.status),
    });
    const list = res.list || [];
    rawMenuList.value = list;
    const menus = list.filter((item: any) => item.type !== 'button');
    const { tree, expandedKeys } = buildTableTreeFromList(menus);
    data.value = tree;
    defaultExpandedRowKeys.value = expandedKeys;
    pagination.itemCount = menus.length || 0;
  } finally {
    loading.value = false;
  }
}

async function fetchMenuTree() {
  const res = await getMenuList({
    page: 1,
    pageSize: 1000,
  });
  const list = res?.list || [];
  menuTreeOptions.value = [
    { id: '0', label: $t('system.menu.form.root') },
    ...buildMenuTreeFromList(list),
  ];
}

function handleSearch() {
  pagination.page = 1;
  fetchMenuList();
}

function resetFilters() {
  filters.name = '';
  filters.status = null;
  pagination.page = 1;
  fetchMenuList();
}

function openCreate() {
  editingId.value = null;
  resetForm();
  fetchMenuTree();
  showModal.value = true;
}

function openEdit(row: any) {
  editingId.value = row.id;
  form.name = row.name ?? '';
  form.path = row.path ?? '';
  form.component = row.component ?? '';
  form.icon = row.icon ?? '';
  form.type = row.type ?? 'menu';
  form.parentId = row.parentId || null;
  if (row.meta && typeof row.meta === 'string') {
    try {
      const parsed = JSON.parse(row.meta);
      form.metaTitle = parsed?.title ?? '';
    } catch {
      form.metaTitle = '';
    }
  } else if (row.meta && typeof row.meta === 'object') {
    form.metaTitle = row.meta?.title ?? '';
  } else {
    form.metaTitle = '';
  }
  form.status = row.status ?? 1;
  form.order = row.order ?? 0;
  fetchMenuTree();
  showModal.value = true;
}

function resetForm() {
  form.name = '';
  form.path = '';
  form.component = '';
  form.icon = '';
  form.type = 'menu';
  form.parentId = null;
  form.metaTitle = '';
  form.status = 1;
  form.order = 0;
  formRef.value?.restoreValidation();
}

function closeModal() {
  showModal.value = false;
}

async function submitForm() {
  if (!formRef.value) {
    return;
  }
  try {
    await formRef.value.validate();
  } catch {
    return;
  }
  saving.value = true;
  try {
    const payload = {
      name: form.name,
      path: form.path,
      component: form.component,
      icon: form.icon,
      type: form.type,
      parentId: form.parentId || 0,
      status: form.status,
      order: form.order,
      ...(form.metaTitle
        ? { meta: JSON.stringify({ title: form.metaTitle }) }
        : {}),
    };
    await (editingId.value
      ? updateMenu(editingId.value, { ...payload, id: editingId.value })
      : createMenu(payload));
    showModal.value = false;
    fetchMenuList();
  } finally {
    saving.value = false;
  }
}

async function handleDelete(row: any) {
  await deleteMenu(row.id);
  fetchMenuList();
}

function openButtonModal(row: any) {
  buttonModal.menuId = String(row.id);
  buttonModal.menuName = getMenuTitle(row);
  buttonEditingId.value = null;
  resetButtonForm();
  buttonList.value =
    rawMenuList.value.filter(
      (item) => String(item.parentId) === buttonModal.menuId && item.type === 'button',
    ) || [];
  buttonModal.show = true;
}

function resetButtonForm() {
  buttonForm.name = '';
  buttonForm.permission = '';
  buttonForm.metaTitle = '';
  buttonForm.status = 1;
  buttonForm.order = 0;
  buttonFormRef.value?.restoreValidation();
}

function openEditButton(item: any) {
  buttonEditingId.value = item.id;
  buttonForm.name = item.name ?? '';
  buttonForm.permission = item.path ?? '';
   buttonForm.metaTitle = getTitleKey(item) ?? '';
  buttonForm.status = item.status ?? 1;
  buttonForm.order = item.order ?? 0;
}

async function submitButtonForm() {
  if (!buttonModal.menuId) return;
  buttonSaving.value = true;
  try {
    const payload = {
      name: buttonForm.name,
      path: buttonForm.permission,
      component: '',
      icon: '',
      type: 'button',
      parentId: buttonModal.menuId,
      status: buttonForm.status,
      order: buttonForm.order,
      ...(buttonForm.metaTitle
        ? { meta: JSON.stringify({ title: buttonForm.metaTitle }) }
        : {}),
    };
    await (buttonEditingId.value
      ? updateMenu(buttonEditingId.value, { ...payload, id: buttonEditingId.value })
      : createMenu(payload));
    await fetchMenuList();
    buttonList.value =
      rawMenuList.value.filter(
        (item) => String(item.parentId) === buttonModal.menuId && item.type === 'button',
      ) || [];
    resetButtonForm();
    buttonEditingId.value = null;
  } finally {
    buttonSaving.value = false;
  }
}

async function handleDeleteButton(item: any) {
  await deleteMenu(item.id);
  await fetchMenuList();
  buttonList.value =
    rawMenuList.value.filter(
      (menu) => String(menu.parentId) === buttonModal.menuId && menu.type === 'button',
    ) || [];
  if (buttonEditingId.value === item.id) {
    resetButtonForm();
    buttonEditingId.value = null;
  }
}

function buildTableTreeFromList(list: any[]) {
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
  function collectExpand(nodes: any[], depth: number) {
    for (const node of nodes) {
      if (depth <= 2) {
        expandedKeys.push(node.id);
      }
      if (node.children?.length) {
        collectExpand(node.children, depth + 1);
      }
    }
  }
  collectExpand(roots, 1);

  return { tree: roots, expandedKeys };
}

onMounted(() => {
  fetchMenuList();
  fetchMenuTree();
});
</script>

<template>
  <div class="menu-page">
    <NCard :title="$t('system.menu.title')" size="small">
      <NForm inline :model="filters" label-placement="left" label-width="auto">
        <NFormItem :label="$t('system.menu.filters.name')">
          <NInput
            v-model:value="filters.name"
            :placeholder="$t('system.menu.filters.searchByName')"
            clearable
          />
        </NFormItem>
        <NFormItem :label="$t('system.menu.filters.status')">
          <NSelect
            v-model:value="filters.status"
            :options="statusOptions"
            :placeholder="$t('system.menu.filters.all')"
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

      <div class="menu-toolbar">
        <NButton type="primary" @click="openCreate">
          {{ $t('system.menu.actions.new') }}
        </NButton>
      </div>

      <NDataTable
        :columns="columns"
        :data="data"
        :pagination="pagination"
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
        <NGrid cols="2" x-gap="16" y-gap="8">
          <NFormItemGi :label="$t('system.menu.form.name')" path="name">
            <NInput
              v-model:value="form.name"
              :placeholder="$t('system.menu.form.namePlaceholder')"
            />
          </NFormItemGi>
          <NFormItemGi :label="$t('system.menu.form.type')" path="type">
            <NSelect v-model:value="form.type" :options="typeOptions" />
          </NFormItemGi>
          <NFormItemGi :label="$t('system.menu.form.path')" path="path">
            <NInput v-model:value="form.path" placeholder="/system/menu" />
          </NFormItemGi>
          <NFormItemGi :label="$t('system.menu.form.titleKey')">
            <NInput
              v-model:value="form.metaTitle"
              :placeholder="$t('system.menu.form.titleKeyPlaceholder')"
            />
          </NFormItemGi>
          <NFormItemGi :label="$t('system.menu.form.component')" path="component">
            <NInput
              v-model:value="form.component"
              placeholder="/system/menu/list"
            />
          </NFormItemGi>
          <NFormItemGi :label="$t('system.menu.form.icon')">
            <NInput
              v-model:value="form.icon"
              :placeholder="$t('system.menu.form.iconPlaceholder')"
            />
          </NFormItemGi>
          <NFormItemGi :label="$t('system.menu.form.parent')">
            <NTreeSelect
              v-model:value="form.parentId"
              :options="menuTreeOptions"
              key-field="id"
              label-field="label"
              children-field="children"
              :placeholder="$t('system.menu.form.root')"
              clearable
            />
          </NFormItemGi>
          <NFormItemGi :label="$t('system.menu.form.order')">
            <NInputNumber v-model:value="form.order" :min="0" />
          </NFormItemGi>
          <NFormItemGi :label="$t('system.menu.form.status')">
            <NSelect v-model:value="form.status" :options="statusOptions" />
          </NFormItemGi>
        </NGrid>
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

    <NModal
      v-model:show="buttonModal.show"
      preset="dialog"
      :title="`${$t('system.menu.actions.manageButtons')} - ${buttonModal.menuName}`"
      :mask-closable="false"
      style="width: 960px"
    >
      <div class="button-toolbar">
        <NButton type="primary" size="small" @click="resetButtonForm">
          {{ $t('common.create') }}
        </NButton>
      </div>
      <NDataTable :columns="buttonColumns" :data="buttonList" :row-key="rowKey" :bordered="false" />
      <div class="button-form">
        <NForm ref="buttonFormRef" :model="buttonForm" label-placement="left" label-width="auto">
          <NGrid cols="2" x-gap="16" y-gap="8">
            <NFormItemGi :label="$t('system.menu.form.name')">
              <NInput
                v-model:value="buttonForm.name"
                :placeholder="$t('system.menu.form.namePlaceholder')"
              />
            </NFormItemGi>
            <NFormItemGi :label="$t('system.menu.form.titleKey')">
              <NInput
                v-model:value="buttonForm.metaTitle"
                :placeholder="$t('system.menu.form.titleKeyPlaceholder')"
              />
            </NFormItemGi>
            <NFormItemGi :label="$t('system.menu.form.path')">
              <NInput
                v-model:value="buttonForm.permission"
                :placeholder="$t('system.menu.form.pathPlaceholder')"
              />
            </NFormItemGi>
            <NFormItemGi :label="$t('system.menu.form.order')">
              <NInputNumber v-model:value="buttonForm.order" :min="0" />
            </NFormItemGi>
            <NFormItemGi :label="$t('system.menu.form.status')">
              <NSelect v-model:value="buttonForm.status" :options="statusOptions" />
            </NFormItemGi>
          </NGrid>
        </NForm>
      </div>
      <template #action>
        <NSpace>
          <NButton @click="buttonModal.show = false">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" :loading="buttonSaving" @click="submitButtonForm">
            {{ $t('common.confirm') }}
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style scoped>
.menu-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.menu-toolbar {
  display: flex;
  justify-content: flex-end;
  margin: 12px 0;
}

.button-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 12px;
}

.button-form {
  margin-top: 12px;
}
</style>
