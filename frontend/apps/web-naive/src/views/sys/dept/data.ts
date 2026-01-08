import type { VbenFormSchema } from '#/adapter/form';

import { z } from '#/adapter/form';
import { $t } from '#/locales';

import { getDeptTree } from '#/api/sys/dept';

export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('system.dept.form.name'),
      rules: 'required',
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
      fieldName: 'parentId',
      label: $t('system.dept.form.parentId'),
    },
    {
      component: 'InputNumber',
      defaultValue: 1,
      fieldName: 'order',
      label: $t('system.dept.form.order'),
      rules: z
        .number({
          required_error: $t('system.dept.validation.orderRequired'),
        })
        .int({ message: $t('system.dept.validation.orderPositive') })
        .min(1, { message: $t('system.dept.validation.orderPositive') }),
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: $t('system.dept.status.enabled'), value: 1 },
          { label: $t('system.dept.status.disabled'), value: 0 },
        ],
      },
      defaultValue: 1,
      fieldName: 'status',
      label: $t('system.dept.form.status'),
    },
  ];
}
