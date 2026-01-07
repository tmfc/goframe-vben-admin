import type { VbenFormSchema } from '#/adapter/form';

import { $t } from '#/locales';

export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('system.dept.form.name'),
      rules: 'required',
    },
    {
      component: 'InputNumber',
      fieldName: 'order',
      label: $t('system.dept.form.order'),
      rules: 'required',
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
