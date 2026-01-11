import type { VbenFormSchema } from '#/adapter/form';

import { $t } from '#/locales';

import { getMenuList } from '#/api/sys/menu';
import { listToTree } from '#/utils/tree';

import MenuIconPicker from './components/MenuIconPicker.vue';
import { getMenuTitle } from './utils';

export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'id',
      hide: true,
      label: 'ID',
    },
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('system.menu.form.name'),
      rules: 'required',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        filterable: true,
        placeholder: $t('system.menu.form.titleKeyPlaceholder'),
        tag: true,
      },
      fieldName: 'metaTitle',
      label: $t('system.menu.form.titleKey'),
    },
    {
      component: 'Input',
      fieldName: 'path',
      label: $t('system.menu.form.path'),
    },
    {
      component: 'Input',
      fieldName: 'component',
      label: $t('system.menu.form.component'),
    },
    {
      component: 'Select',
      componentProps: {
        options: [
          { label: $t('system.menu.types.menu'), value: 'menu' },
          { label: $t('system.menu.types.catalog'), value: 'catalog' },
          { label: $t('system.menu.types.link'), value: 'link' },
          { label: $t('system.menu.types.embedded'), value: 'embedded' },
        ],
      },
      fieldName: 'type',
      label: $t('system.menu.form.type'),
      rules: 'required',
    },
    {
      component: MenuIconPicker as any,
      fieldName: 'icon',
      label: $t('system.menu.form.icon'),
    },
    {
            component: 'ApiTreeSelect',
            componentProps: {
              api: async () => {
                const res = await getMenuList();
                const list = (res.list || []).map((item) => ({
                  ...item,
                  name: getMenuTitle(item),
                }));
                const tree = listToTree(list);
                return [
                  {
                    id: '0',
                    name: $t('system.menu.form.root'),
                  },
                  ...tree,
                ];
              },        childrenField: 'children',
        labelField: 'name',
        placeholder: $t('system.menu.form.root'),
        valueField: 'id',
      },
      fieldName: 'parentId',
      label: $t('system.menu.form.parent'),
    },
    {
      component: 'InputNumber',
      componentProps: {
        min: 0,
      },
      fieldName: 'order',
      label: $t('system.menu.form.order'),
    },
    {
      component: 'Select',
      componentProps: {
        options: [
          { label: $t('system.menu.status.enabled'), value: 1 },
          { label: $t('system.menu.status.disabled'), value: 0 },
        ],
      },
      fieldName: 'status',
      label: $t('system.menu.form.status'),
    },
    {
      component: 'Input',
      fieldName: 'permissionCode',
      label: $t('system.menu.form.permissionCode'),
    },
    {
      component: 'Checkbox',
      fieldName: 'autoGeneratePermission',
      label: $t('system.menu.form.autoGeneratePermission'),
    },
  ];
}
