import type { VbenFormSchema } from '#/adapter/form';

import { $t } from '#/locales';

import { getDeptTree } from '#/api/sys/dept';

export function useFormSchema(isUpdate?: () => boolean): VbenFormSchema[] {
  const shouldShowPassword = () => !(isUpdate?.() ?? false);
  return [
    {
      component: 'Input',
      fieldName: 'id',
      hide: true,
      label: 'ID',
    },
    {
      component: 'Input',
      fieldName: 'username',
      label: $t('system.user.form.username'),
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'realName',
      label: $t('system.user.form.realName'),
    },
    {
      component: 'Input',
      componentProps: {
        type: 'password',
        showPasswordOn: 'click',
      },
      fieldName: 'password',
      label: $t('system.user.form.password'),
      dependencies: {
        if: () => shouldShowPassword(),
        triggerFields: [],
      },
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: $t('system.user.form.rolesPlaceholder'),
      },
      fieldName: 'roles',
      label: $t('system.user.form.roles'),
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: $t('system.user.form.homePathPlaceholder'),
      },
      fieldName: 'homePath',
      label: $t('system.user.form.homePath'),
    },
    {
      component: 'Input',
      componentProps: {
        placeholder: $t('system.user.form.avatarPlaceholder'),
      },
      fieldName: 'avatar',
      label: $t('system.user.form.avatar'),
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
      label: $t('system.user.form.deptId'),
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: $t('system.user.status.enabled'), value: 1 },
          { label: $t('system.user.status.disabled'), value: 0 },
        ],
      },
      defaultValue: 1,
      fieldName: 'status',
      label: $t('system.user.form.status'),
    },
  ];
}
