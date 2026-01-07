import type { VbenFormSchema } from '#/adapter/form';

export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'name',
      label: 'Department Name',
      rules: 'required',
    },
    {
      component: 'InputNumber',
      fieldName: 'order',
      label: 'Order',
      rules: 'required',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: 'Enabled', value: 1 },
          { label: 'Disabled', value: 0 },
        ],
      },
      defaultValue: 1,
      fieldName: 'status',
      label: 'Status',
    },
  ];
}
