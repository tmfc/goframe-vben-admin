import type { VbenFormSchema } from '#/adapter/form';

import { $t } from '#/locales';

import { getDeptTree } from '#/api/sys/dept';
import { getRoleList } from '#/api/sys/role';

export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('system.role.form.name'),
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'description',
      label: $t('system.role.form.description'),
    },
    {
      component: 'ApiSelect',
      componentProps: {
        api: getRoleList,
        resultField: 'items',
        labelField: 'name',
        valueField: 'id',
        params: {
          page: 1,
          pageSize: 5000,
        },
        clearable: true,
      },
      fieldName: 'parentId',
      label: $t('system.role.form.parentId'),
    },
    {
      component: 'ApiTreeSelect',
      componentProps: {
        api: getDeptTree,
        resultField: 'list',
        labelField: 'name',
        valueField: 'id',
        childrenField: 'children',
        clearable: true,
      },
      fieldName: 'deptId',
      label: $t('system.role.form.deptId'),
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: $t('system.role.status.enabled'), value: 1 },
          { label: $t('system.role.status.disabled'), value: 0 },
        ],
      },
      defaultValue: 1,
      fieldName: 'status',
      label: $t('system.role.form.status'),
    },
  ];
}
