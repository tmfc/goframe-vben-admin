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
  NTag,
  useDialog,
} from 'naive-ui';

import { $t } from '#/locales';

import { getDeptTree } from '#/api/sys/dept';
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
      title: $t('system.role.columns.name'),
      key: 'name',
    },
    {
      title: $t('system.role.columns.parent'),
      key: 'parentId',
      render(row: Role) {
        const parentId = row?.parentId;
        if (!parentId) {
          return '';
        }
        return roleNameMap.value.get(String(parentId)) || String(parentId);
      },
    },
    {
      title: $t('system.role.columns.description'),
      key: 'description',
    },
    {
      title: $t('system.role.columns.dept'),
      key: 'deptId',
      render(row: Role) {
        const deptId = row?.deptId;
        if (!deptId) {
          return '';
        }
        return deptNameMap.value.get(String(deptId)) || String(deptId);
      },
    },
    {
      title: $t('system.role.columns.status'),
      key: 'status',
      render(row: Role) {
        const isActive = row.status === 1;
        return h(
          NTag,
          { type: isActive ? 'success' : 'default', size: 'small' },
          {
            default: () =>
              isActive
                ? $t('system.role.status.enabled')
                : $t('system.role.status.disabled'),
          },
        );
      },
    },
    {
      title: $t('system.role.columns.actions'),
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

const data = ref<Role[]>([]);
const deptNameMap = ref(new Map<string, string>());
const roleNameMap = ref(new Map<string, string>());
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
    title: $t('system.role.dialog.deleteTitle'),
    content: $t('system.role.dialog.deleteConfirm'),
    positiveText: $t('common.confirm'),
    negativeText: $t('common.cancel'),
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

async function fetchRoleNameMap() {
  try {
    const response = await getRoleList({ page: 1, pageSize: 5000 });
    roleNameMap.value = new Map(
      (response.items || []).map((item: Role) => [
        String(item.id),
        item.name ?? '',
      ]),
    );
  } catch (error) {
    console.error(error);
  }
}

const columns = createColumns({
  edit: handleEdit,
  del: handleDelete,
});

onMounted(() => {
  getDeptTree()
    .then((res) => {
      deptNameMap.value = flattenDeptTree(res?.list || []);
    })
    .catch((error) => {
      console.error(error);
    });
  fetchRoleNameMap();
  fetchData();
});
</script>

<template>
  <div class="role-page">
    <NCard :title="$t('system.role.title')" size="small">
      <div class="role-toolbar">
        <NForm
          inline
          :model="formValue"
          label-placement="left"
          label-width="auto"
          @submit.prevent="handleSearch"
        >
          <NFormItem :label="$t('system.role.filters.name')">
            <NInput
              v-model:value="formValue.name"
              :placeholder="$t('system.role.filters.searchByName')"
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
          {{ $t('system.role.actions.create') }}
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
