import type { VbenFormSchema } from '#/adapter/form';

export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'name',
      label: 'Role Name',
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'description',
      label: 'Description',
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
