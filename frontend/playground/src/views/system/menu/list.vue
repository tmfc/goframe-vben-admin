<script lang="ts" setup>
import type {
  OnActionClickParams,
  VxeTableGridOptions,
} from '#/adapter/vxe-table';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { IconifyIcon, Plus, Zap } from '@vben/icons';
import { $t } from '@vben/locales';

import { MenuBadge } from '@vben-core/menu-ui';

import { Button, message } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { createMenu, deleteMenu, getMenuList, SystemMenuApi } from '#/api/system/menu';

import { useColumns } from './data';
import Form from './modules/form.vue';

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  connectedComponent: Form,
  destroyOnClose: true,
});

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions: {
    columns: useColumns(onActionClick),
    height: 'auto',
    keepSource: true,
    checkboxConfig: {
      highlight: true,
      range: true,
    },
    pagerConfig: {
      enabled: false,
    },
    proxyConfig: {
      ajax: {
        query: async (_params) => {
          return await getMenuList();
        },
      },
    },
    rowConfig: {
      keyField: 'id',
    },
    sortConfig: {
      defaultSort: {
        field: 'order',
        order: 'asc',
      },
      remote: false,
      trigger: 'cell',
    },
    toolbarConfig: {
      custom: true,
      export: false,
      refresh: true,
      zoom: true,
    },
    treeConfig: {
      parentField: 'pid',
      rowField: 'id',
      transform: false,
    },
  } as VxeTableGridOptions,
});

function onActionClick({
  code,
  row,
}: OnActionClickParams<SystemMenuApi.SystemMenu>) {
  switch (code) {
    case 'append': {
      onAppend(row);
      break;
    }
    case 'delete': {
      onDelete(row);
      break;
    }
    case 'edit': {
      onEdit(row);
      break;
    }
    default: {
      break;
    }
  }
}

function onRefresh() {
  gridApi.query();
}
function onEdit(row: SystemMenuApi.SystemMenu) {
  formDrawerApi.setData(row).open();
}
function onCreate() {
  formDrawerApi.setData({}).open();
}
function onAppend(row: SystemMenuApi.SystemMenu) {
  formDrawerApi.setData({ pid: row.id }).open();
}

function onDelete(row: SystemMenuApi.SystemMenu) {
  const hideLoading = message.loading({
    content: $t('ui.actionMessage.deleting', [row.name]),
    duration: 0,
    key: 'action_process_msg',
  });
  deleteMenu(row.id)
    .then(() => {
      message.success({
        content: $t('ui.actionMessage.deleteSuccess', [row.name]),
        key: 'action_process_msg',
      });
      onRefresh();
    })
    .catch(() => {
      hideLoading();
    });
}

async function onGenerateButtons() {
  const selectedRows = gridApi.getCheckboxRecords();
  if (selectedRows.length === 0) {
    message.warning('请先选择一个菜单项');
    return;
  }

  const parentMenu = selectedRows[0] as SystemMenuApi.SystemMenu;

  // 只能为 menu 类型生成按钮
  if (parentMenu.type !== 'menu') {
    message.warning('只能为菜单类型生成按钮');
    return;
  }

  const hideLoading = message.loading({
    content: '正在生成按钮...',
    duration: 0,
    key: 'generate_buttons',
  });

  try {
    // 生成三个按钮: 新增、编辑、删除
    const buttons = [
      {
        name: `${parentMenu.name}Create`,
        type: 'button',
        status: 1,
        authCode: `Button:${parentMenu.name}:Create`,
        pid: parentMenu.id,
        meta: {
          title: 'common.create',
        },
        order: 0,
      },
      {
        name: `${parentMenu.name}Edit`,
        type: 'button',
        status: 1,
        authCode: `Button:${parentMenu.name}:Edit`,
        pid: parentMenu.id,
        meta: {
          title: 'common.edit',
        },
        order: 1,
      },
      {
        name: `${parentMenu.name}Delete`,
        type: 'button',
        status: 1,
        authCode: `Button:${parentMenu.name}:Delete`,
        pid: parentMenu.id,
        meta: {
          title: 'common.delete',
        },
        order: 2,
      },
    ];

    // 依次创建按钮
    for (const button of buttons) {
      await createMenu(button);
    }

    message.success({
      content: '按钮生成成功',
      key: 'generate_buttons',
    });
    onRefresh();
  } catch (error) {
    message.error({
      content: '按钮生成失败',
      key: 'generate_buttons',
    });
  } finally {
    hideLoading();
  }
}
</script>
<template>
  <Page auto-content-height>
    <FormDrawer @success="onRefresh" />
    <Grid>
      <template #toolbar-tools>
        <Button type="primary" @click="onCreate">
          <Plus class="size-5" />
          {{ $t('ui.actionTitle.create', [$t('system.menu.name')]) }}
        </Button>
        <Button @click="onGenerateButtons">
          <Zap class="size-5" />
          生成按钮
        </Button>
      </template>
      <template #title="{ row }">
        <div class="flex w-full items-center gap-1">
          <div class="size-5 flex-shrink-0">
            <IconifyIcon
              v-if="row.type === 'button'"
              icon="carbon:security"
              class="size-full"
            />
            <IconifyIcon
              v-else-if="row.meta?.icon"
              :icon="row.meta?.icon || 'carbon:circle-dash'"
              class="size-full"
            />
          </div>
          <span class="flex-auto">{{ $t(row.meta?.title) }}</span>
          <div class="items-center justify-end"></div>
        </div>
        <MenuBadge
          v-if="row.meta?.badgeType"
          class="menu-badge"
          :badge="row.meta.badge"
          :badge-type="row.meta.badgeType"
          :badge-variants="row.meta.badgeVariants"
        />
      </template>
    </Grid>
  </Page>
</template>
<style lang="scss" scoped>
.menu-badge {
  top: 50%;
  right: 0;
  transform: translateY(-50%);

  & > :deep(div) {
    padding-top: 0;
    padding-bottom: 0;
  }
}
</style>
