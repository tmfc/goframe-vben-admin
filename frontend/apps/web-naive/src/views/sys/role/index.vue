<script lang="ts" setup>
import type { PaginationProps } from 'naive-ui';
import type { Role } from '#/api/sys/role';
import type { Permission } from '#/api/sys/permission';

import { h, onMounted, ref } from 'vue';

import {
  NButton,
  NCard,
  NDataTable,
  NEmpty,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NSpace,
  NSpin,
  NTag,
  NTree,
  useDialog,
} from 'naive-ui';

import { $t } from '#/locales';

import { getDeptTree } from '#/api/sys/dept';
import { getPermissionList } from '#/api/sys/permission';
import {
  assignRolePermissions,
  deleteRole,
  getRoleList,
  getRolePermissions,
} from '#/api/sys/role';
import RoleFormModal from './modules/form.vue';

const dialog = useDialog();

function createColumns({
  edit,
  del,
  permissions,
}: {
  del: (id: string) => void;
  edit: (record: Role) => void;
  permissions: (record: Role) => void;
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
                  onClick: () => edit(row),
                },
                { default: () => $t('common.edit') },
              ),
              h(
                NButton,
                {
                  size: 'small',
                  tertiary: true,
                  onClick: () => permissions(row),
                },
                { default: () => $t('system.role.actions.permissions') },
              ),
              h(
                NButton,
                {
                  size: 'small',
                  type: 'error',
                  secondary: true,
                  onClick: () => del(row.id),
                },
                { default: () => $t('common.delete') },
              ),
            ],
          },
        );
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
const permissionModal = ref({
  show: false,
  roleId: '',
  roleName: '',
});
const permissionTree = ref<Array<{ key: string; label: string; children?: any[] }>>(
  [],
);
const permissionCheckedKeys = ref<string[]>([]);
const permissionExpandedKeys = ref<string[]>([]);
const permissionLoading = ref(false);
const permissionSaving = ref(false);

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

function buildPermissionTree(list: Permission[] = []) {
  const nodeMap = new Map<string, { key: string; label: string; children: any[] }>();
  const roots: Array<{ key: string; label: string; children: any[] }> = [];

  for (const item of list || []) {
    if (!item?.id) continue;
    const key = String(item.id);
    nodeMap.set(key, {
      key,
      label: item.name ?? key,
      children: [],
    });
  }

  for (const item of list || []) {
    if (!item?.id) continue;
    const key = String(item.id);
    const parentId = item.parentId ? String(item.parentId) : '';
    const node = nodeMap.get(key);
    if (!node) continue;
    if (!parentId || parentId === '0' || !nodeMap.has(parentId)) {
      roots.push(node);
    } else {
      nodeMap.get(parentId)?.children.push(node);
    }
  }

  const expandedKeys: string[] = [];
  const walk = (nodes: any[], depth: number) => {
    for (const node of nodes) {
      if (depth <= 2) {
        expandedKeys.push(node.key);
      }
      if (node.children?.length) {
        walk(node.children, depth + 1);
      }
    }
  };
  walk(roots, 1);

  return { tree: roots, expandedKeys };
}

async function fetchPermissionList() {
  const permissionRes = await getPermissionList({ page: 1, pageSize: 1000 });
  const list = permissionRes?.items || permissionRes?.list || [];
  const { tree, expandedKeys } = buildPermissionTree(list);
  permissionTree.value = tree;
  permissionExpandedKeys.value = expandedKeys;
}

async function openPermissionModal(record: Role) {
  permissionModal.value = {
    show: true,
    roleId: String(record.id),
    roleName: record.name ?? '',
  };
  permissionLoading.value = true;
  try {
    const [_, rolePermRes] = await Promise.all([
      fetchPermissionList(),
      getRolePermissions(String(record.id)),
    ]);
    permissionCheckedKeys.value = (rolePermRes?.permissionIds || []).map(String);
  } catch (error) {
    console.error(error);
  } finally {
    permissionLoading.value = false;
  }
}

function closePermissionModal() {
  permissionModal.value.show = false;
}

async function submitPermissionChange() {
  if (!permissionModal.value.roleId) return;
  permissionSaving.value = true;
  try {
    const permissionIds = permissionCheckedKeys.value
      .map((key) => Number(key))
      .filter((value) => Number.isFinite(value));
    await assignRolePermissions(permissionModal.value.roleId, { permissionIds });
    closePermissionModal();
  } catch (error) {
    console.error(error);
  } finally {
    permissionSaving.value = false;
  }
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
  permissions: openPermissionModal,
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
    <NModal
      v-model:show="permissionModal.show"
      preset="dialog"
      :title="`${$t('system.role.permissions.title')} - ${permissionModal.roleName}`"
      :mask-closable="false"
      style="width: 720px"
    >
      <NSpin :show="permissionLoading">
        <NTree
          v-if="permissionTree.length"
          v-model:checked-keys="permissionCheckedKeys"
          :data="permissionTree"
          :default-expanded-keys="permissionExpandedKeys"
          checkable
          block-line
        />
        <NEmpty v-else :description="$t('system.role.permissions.empty')" />
      </NSpin>
      <template #action>
        <NSpace>
          <NButton @click="closePermissionModal">
            {{ $t('common.cancel') }}
          </NButton>
          <NButton type="primary" :loading="permissionSaving" @click="submitPermissionChange">
            {{ $t('common.confirm') }}
          </NButton>
        </NSpace>
      </template>
    </NModal>
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
