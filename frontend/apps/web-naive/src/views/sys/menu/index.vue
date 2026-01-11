<script lang="ts" setup>
import type { DataTableColumns, FormInst } from 'naive-ui';

import { h, onMounted, reactive, ref, watch } from 'vue';
import { useI18n } from '@vben/locales';

import {
  NButton,
  NCard,
  NCheckbox,
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
  useMessage,
} from 'naive-ui';

import MenuFormModal from './modules/form.vue';
import { generatePermissionCode, useTitleOptions } from './utils';

import {
  createMenu,
  deleteMenu,
  generateButtons,
  getMenuList,
  updateMenu,
  type MenuItem,
} from '#/api/sys/menu';
import { collectExpandedKeys, listToTree, sortTree } from '#/utils/tree';
import { Grip, IconifyIcon } from '@vben/icons';
import zhSystem from '#/locales/langs/zh-CN/system.json';
import enSystem from '#/locales/langs/en-US/system.json';
import { $t } from '#/locales';

const loading = ref(false);
const generating = ref(false);
const showModal = ref(false);
const editingRecord = ref<MenuItem | null>(null);
const message = useMessage();

const filters = reactive({
  name: '',
  status: null as null | number,
});

const { locale } = useI18n();
const localeSystemMap: Record<string, any> = {
  'zh-CN': zhSystem,
  'en-US': enSystem,
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
const rawMenuList = ref<MenuItem[]>([]);
const titleKeyOptions = useTitleOptions();

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
      return h(
        'div',
        {
          class: 'menu-name-cell',
          style: {
            display: 'inline-flex',
            alignItems: 'center',
            gap: '8px',
            whiteSpace: 'nowrap',
          },
        },
        {
          default: () => [
            row.icon
              ? h(IconifyIcon, {
                  icon: row.icon,
                  class: 'menu-name-cell__icon',
                  style: { width: '16px', height: '16px', flex: '0 0 auto' },
                })
              : h(Grip, {
                  class: 'menu-name-cell__icon',
                  style: { width: '16px', height: '16px', flex: '0 0 auto' },
                }),
            h(
              'span',
              { class: 'menu-name-cell__text', style: { whiteSpace: 'nowrap' } },
              getMenuTitle(row),
            ),
          ],
        },
      );
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
  {
    title: $t('system.menu.columns.permissionCode'),
    key: 'permissionCode',
    render(row) {
      return row.permissionCode || '-';
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
    title: $t('system.menu.columns.permissionCode'),
    key: 'permissionCode',
    render(row) {
      return row.permissionCode || '-';
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
  titleKey: '',
  metaTitle: '',
  status: 1,
  order: 0,
  permissionCode: '',
  autoGeneratePermission: true,
});

async function fetchMenuList() {
  try {
    loading.value = true;
    const res = await getMenuList({
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
  editingRecord.value = null;
  showModal.value = true;
}

function openEdit(row: MenuItem) {
  editingRecord.value = row;
  showModal.value = true;
}

async function handleDelete(row: MenuItem) {
  await deleteMenu(row.id);
  fetchMenuList();
}

function openButtonModal(row: any) {
  buttonModal.menuId = String(row.id);
  buttonModal.menuName = getMenuTitle(row);
  buttonEditingId.value = null;
  resetButtonForm();
  buttonList.value =
    rawMenuList.value
      .filter((item) => String(item.parentId) === buttonModal.menuId && item.type === 'button')
      .sort((a, b) => (a.order ?? 0) - (b.order ?? 0)) || [];
  buttonModal.show = true;
}

async function generateDefaultButtons() {
  if (!buttonModal.menuId) return;
  generating.value = true;
  try {
    await generateButtons(buttonModal.menuId);
    message.success($t('system.menu.actions.generateSuccess'));
    await fetchMenuList();
    buttonList.value =
      rawMenuList.value
        .filter((item) => String(item.parentId) === buttonModal.menuId && item.type === 'button')
        .sort((a, b) => (a.order ?? 0) - (b.order ?? 0)) || [];
  } catch (error) {
    message.error($t('system.menu.actions.generateFailed'));
  } finally {
    generating.value = false;
  }
}

function resetButtonForm() {
  buttonForm.name = '';
  buttonForm.permission = '';
  buttonForm.metaTitle = '';
  buttonForm.status = 1;
  buttonForm.order = 0;
  buttonForm.permissionCode = '';
  buttonForm.autoGeneratePermission = true;
  buttonFormRef.value?.restoreValidation();
}

function openEditButton(item: any) {
  buttonEditingId.value = item.id;
  buttonForm.name = item.name ?? '';
  buttonForm.permission = item.path ?? '';
   buttonForm.metaTitle = getTitleKey(item) ?? '';
  buttonForm.status = item.status ?? 1;
  buttonForm.order = item.order ?? 0;
  buttonForm.permissionCode = item.permissionCode ?? '';
  buttonForm.autoGeneratePermission = !item.permissionCode;
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
      permissionCode: buttonForm.permissionCode || undefined,
      ...(buttonForm.metaTitle
        ? { meta: JSON.stringify({ title: buttonForm.metaTitle }) }
        : {}),
    };
    await (buttonEditingId.value
      ? updateMenu(buttonEditingId.value, { ...payload, id: buttonEditingId.value })
      : createMenu(payload));
    await fetchMenuList();
    buttonList.value =
      rawMenuList.value
        .filter((item) => String(item.parentId) === buttonModal.menuId && item.type === 'button')
        .sort((a, b) => (a.order ?? 0) - (b.order ?? 0)) || [];
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
    rawMenuList.value
      .filter((menu) => String(menu.parentId) === buttonModal.menuId && menu.type === 'button')
      .sort((a, b) => (a.order ?? 0) - (b.order ?? 0)) || [];
  if (buttonEditingId.value === item.id) {
    resetButtonForm();
    buttonEditingId.value = null;
  }
}

function rebuildTitleKeyOptionsFromLocale() {
  const set = new Set<string>();
  function collect(obj: any, prefix = '') {
    if (!obj || typeof obj !== 'object') return;
    Object.entries(obj).forEach(([k, v]) => {
      const keyPath = prefix ? `${prefix}${k}` : k;
      if (typeof v === 'string') {
        set.add(keyPath);
        return;
      }
      if (v && typeof v === 'object') {
        collect(v, `${keyPath}.`);
      }
    });
  }
  const currentLocale = locale.value;
  const currentSystem = localeSystemMap[currentLocale] || localeSystemMap['zh-CN'];
  collect(currentSystem?.common, 'common.');
  collect(currentSystem?.menu, 'system.menu.');
  collect(currentSystem?.permission, 'system.permission.');
  titleKeyOptions.value = Array.from(set).map((key) => ({
    label: `${$t(key)} (${key})`,
    value: key,
  }));
}

watch(
  () => [buttonForm.name, buttonForm.autoGeneratePermission, buttonModal.menuId],
  ([name, autoGenerate]) => {
    if (autoGenerate && name) {
      buttonForm.permissionCode = generatePermissionCode(
        String(name),
        'button',
        buttonModal.menuId,
        rawMenuList.value,
      );
    }
  },
);

function buildTableTreeFromList(list: any[]) {
  const tree = listToTree(list);
  sortTree(tree, (a, b) => (a.order ?? 0) - (b.order ?? 0));
  const expandedKeys = collectExpandedKeys(tree, 2);
  return { tree, expandedKeys };
}

onMounted(() => {
  fetchMenuList();
  rebuildTitleKeyOptionsFromLocale();
});

watch(
  () => locale.value,
  () => {
    rebuildTitleKeyOptionsFromLocale();
  },
);
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

    <MenuFormModal
      v-model:show="showModal"
      :record="editingRecord"
      @success="fetchMenuList"
    />

    <NModal
      v-model:show="buttonModal.show"
      preset="dialog"
      :title="`${$t('system.menu.actions.manageButtons')} - ${buttonModal.menuName}`"
      :mask-closable="false"
      style="width: 960px"
    >
      <div class="button-toolbar">
        <NSpace>
          <NButton size="small" :loading="generating" @click="generateDefaultButtons">
            {{ $t('system.menu.actions.generateButtons') }}
          </NButton>
          <NButton type="primary" size="small" @click="resetButtonForm">
            {{ $t('common.create') }}
          </NButton>
        </NSpace>
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
              <NSelect
                v-model:value="buttonForm.metaTitle"
                :options="titleKeyOptions"
                filterable
                clearable
                tag
                :placeholder="$t('system.menu.form.titleKeyPlaceholder')"
              />
            </NFormItemGi>
            <NFormItemGi :label="$t('system.menu.form.path')">
              <NInput
                v-model:value="buttonForm.permission"
                :placeholder="$t('system.menu.form.pathPlaceholder')"
              />
            </NFormItemGi>
            <NFormItemGi :label="$t('system.menu.form.permissionCode')">
              <NInput
                v-model:value="buttonForm.permissionCode"
                :placeholder="$t('system.menu.form.permissionCodePlaceholder')"
                :disabled="buttonForm.autoGeneratePermission"
              />
            </NFormItemGi>
            <NFormItemGi :label="$t('system.menu.form.order')">
              <NInputNumber v-model:value="buttonForm.order" :min="0" />
            </NFormItemGi>
            <NFormItemGi :label="$t('system.menu.form.status')">
              <NSelect v-model:value="buttonForm.status" :options="statusOptions" />
            </NFormItemGi>
            <NFormItemGi :label="''">
              <NCheckbox v-model:checked="buttonForm.autoGeneratePermission">
                {{ $t('system.menu.form.autoGeneratePermission') }}
              </NCheckbox>
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

.menu-name-cell {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.menu-name-cell__icon {
  width: 16px;
  height: 16px;
  flex: 0 0 auto;
  color: var(--n-text-color-2);
}

.menu-name-cell__text {
  line-height: 1.2;
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