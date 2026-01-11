import type { VbenFormSchema } from '#/adapter/form';

import { $t } from '#/locales';

import { getPermissionList } from '#/api/sys/permission';
import { listToTree } from '#/utils/tree';

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
      label: $t('system.permission.form.name'),
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'description',
      label: $t('system.permission.form.description'),
    },
    {
      component: 'ApiTreeSelect',
      componentProps: {
        api: async () => {
          const res = await getPermissionList();
          const list = res.items || [];
          const tree = listToTree(list);
          return [
            {
              id: '0',
              name: $t('system.permission.form.root'),
            },
            ...tree,
          ];
        },
        childrenField: 'children',
        labelField: 'name',
        placeholder: $t('system.permission.form.root'),
        valueField: 'id',
      },
      fieldName: 'parentId',
      label: $t('system.permission.form.parent'),
    },
    {
      component: 'Select',
      componentProps: {
        options: [
          { label: $t('system.permission.status.enabled'), value: 1 },
          { label: $t('system.permission.status.disabled'), value: 0 },
        ],
      },
      fieldName: 'status',
      label: $t('system.permission.form.status'),
    },
  ];
}
