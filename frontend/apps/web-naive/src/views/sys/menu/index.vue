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
  name: [{ required: true, message: 'Name is required', trigger: 'blur' }],
  type: [{ required: true, message: 'Type is required', trigger: 'change' }],
};

const statusOptions = [
  { label: 'Enabled', value: 1 },
  { label: 'Disabled', value: 0 },
];

const typeOptions = [
  { label: 'Menu', value: 'menu' },
  { label: 'Catalog', value: 'catalog' },
  { label: 'Link', value: 'link' },
  { label: 'Embedded', value: 'embedded' },
  { label: 'Button', value: 'button' },
];

const data = ref<any[]>([]);

const parentOptions = computed(() =>
  data.value.map((item) => ({
    label: item.name,
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
  { title: 'Name', key: 'name' },
  { title: 'Path', key: 'path' },
  { title: 'Component', key: 'component' },
  { title: 'Type', key: 'type' },
  { title: 'Order', key: 'order' },
  {
    title: 'Status',
    key: 'status',
    render(row) {
      const isActive = row.status === 1;
      return h(
        NTag,
        { type: isActive ? 'success' : 'default', size: 'small' },
        { default: () => (isActive ? 'Enabled' : 'Disabled') },
      );
    },
  },
  {
    title: 'Action',
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
              { default: () => 'Edit' },
            ),
            h(
              NPopconfirm,
              { onPositiveClick: () => handleDelete(row) },
              {
                trigger: () =>
                  h(
                    NButton,
                    { size: 'small', type: 'error', secondary: true },
                    { default: () => 'Delete' },
                  ),
                default: () => 'Delete this menu?',
              },
            ),
          ],
        },
      );
    },
  },
]);

const modalTitle = computed(() =>
  editingId.value ? 'Edit Menu' : 'Create Menu',
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
    <NCard title="Menu Management" size="small">
      <NForm inline :model="filters" label-placement="left" label-width="auto">
        <NFormItem label="Name">
          <NInput
            v-model:value="filters.name"
            placeholder="Search by name"
            clearable
          />
        </NFormItem>
        <NFormItem label="Status">
          <NSelect
            v-model:value="filters.status"
            :options="statusOptions"
            placeholder="All"
            clearable
          />
        </NFormItem>
        <NFormItem>
          <NSpace>
            <NButton type="primary" @click="handleSearch">Search</NButton>
            <NButton @click="resetFilters">Reset</NButton>
          </NSpace>
        </NFormItem>
      </NForm>

      <div class="menu-toolbar">
        <NButton type="primary" @click="openCreate">New Menu</NButton>
      </div>

      <NDataTable
        :columns="columns"
        :data="data"
        :pagination="pagination"
        :loading="loading"
        :bordered="false"
        row-key="id"
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
          <NFormItemGi label="Name" path="name">
            <NInput v-model:value="form.name" placeholder="Menu name" />
          </NFormItemGi>
          <NFormItemGi label="Type" path="type">
            <NSelect v-model:value="form.type" :options="typeOptions" />
          </NFormItemGi>
          <NFormItemGi label="Path" path="path">
            <NInput v-model:value="form.path" placeholder="/system/menu" />
          </NFormItemGi>
          <NFormItemGi label="Component" path="component">
            <NInput
              v-model:value="form.component"
              placeholder="/system/menu/list"
            />
          </NFormItemGi>
          <NFormItemGi label="Icon">
            <NInput v-model:value="form.icon" placeholder="icon name" />
          </NFormItemGi>
          <NFormItemGi label="Parent">
            <NSelect
              v-model:value="form.parentId"
              :options="parentOptions"
              placeholder="Root"
              clearable
            />
          </NFormItemGi>
          <NFormItemGi label="Order">
            <NInputNumber v-model:value="form.order" :min="0" />
          </NFormItemGi>
          <NFormItemGi label="Status">
            <NSelect v-model:value="form.status" :options="statusOptions" />
          </NFormItemGi>
        </NGrid>
      </NForm>
      <template #action>
        <NSpace>
          <NButton @click="closeModal">Cancel</NButton>
          <NButton type="primary" :loading="saving" @click="submitForm">
            Save
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
