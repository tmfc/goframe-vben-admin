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
  { label: $t('system.menu.types.button'), value: 'button' },
];

const data = ref<any[]>([]);

function getMenuTitle(row: any) {
  if (!row) {
    return '';
  }
  const rawMeta = row.meta;
  let meta: Record<string, any> | null = null;
  if (typeof rawMeta === 'string') {
    try {
      meta = JSON.parse(rawMeta);
    } catch {
      meta = null;
    }
  } else if (rawMeta && typeof rawMeta === 'object') {
    meta = rawMeta;
  }
  const titleKey = meta?.title;
  if (titleKey) {
    return $t(titleKey);
  }
  return row.name ?? '';
}

const parentOptions = computed(() =>
  data.value.map((item) => ({
    label: getMenuTitle(item),
    value: item.id,
  })),
);

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

const rowKey = (row: any) => row.id;

const modalTitle = computed(() =>
  editingId.value
    ? $t('system.menu.actions.edit')
    : $t('system.menu.actions.create'),
);

async function fetchMenuList() {
  try {
    loading.value = true;
    const res = await getMenuList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      name: filters.name || undefined,
      status: filters.status === null ? undefined : String(filters.status),
    });
    data.value = res.list || [];
    pagination.itemCount = res.total || 0;
  } finally {
    loading.value = false;
  }
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
  form.status = row.status ?? 1;
  form.order = row.order ?? 0;
  showModal.value = true;
}

function resetForm() {
  form.name = '';
  form.path = '';
  form.component = '';
  form.icon = '';
  form.type = 'menu';
  form.parentId = null;
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

onMounted(() => {
  fetchMenuList();
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
            <NSelect
              v-model:value="form.parentId"
              :options="parentOptions"
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
</style>
