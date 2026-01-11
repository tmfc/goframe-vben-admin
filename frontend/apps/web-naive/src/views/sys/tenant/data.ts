import type { VbenFormSchema } from '#/adapter/form';
import { $t } from '#/locales';

export const useFormSchema = (): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: {
      placeholder: $t('system.tenant.form.name'),
    },
    fieldName: 'name',
    label: $t('system.tenant.form.name'),
    rules: 'required',
  },
  {
    component: 'Input',
    componentProps: {
      placeholder: $t('system.tenant.form.code'),
    },
    fieldName: 'code',
    label: $t('system.tenant.form.code'),
    rules: 'required',
  },
  {
    component: 'Input',
    componentProps: {
      placeholder: $t('system.tenant.form.contactName'),
    },
    fieldName: 'contactName',
    label: $t('system.tenant.form.contactName'),
  },
  {
    component: 'Input',
    componentProps: {
      placeholder: $t('system.tenant.form.contactMobile'),
    },
    fieldName: 'contactMobile',
    label: $t('system.tenant.form.contactMobile'),
  },
  {
    component: 'DatePicker',
    componentProps: {
      type: 'datetime',
      valueFormat: 'yyyy-MM-dd HH:mm:ss',
    },
    fieldName: 'startAt',
    label: $t('system.tenant.form.startAt'),
  },
  {
    component: 'DatePicker',
    componentProps: {
      type: 'datetime',
      valueFormat: 'yyyy-MM-dd HH:mm:ss',
    },
    fieldName: 'expireAt',
    label: $t('system.tenant.form.expireAt'),
  },
  {
    component: 'Input',
    componentProps: {
      placeholder: $t('system.tenant.form.packageVersion'),
    },
    fieldName: 'packageVersion',
    label: $t('system.tenant.form.packageVersion'),
  },
  {
    component: 'Input',
    componentProps: {
      placeholder: $t('system.tenant.form.domain'),
    },
    fieldName: 'domain',
    label: $t('system.tenant.form.domain'),
  },
  {
    component: 'RadioGroup',
    componentProps: {
      options: [
        { label: $t('system.tenant.status.enabled'), value: 1 },
        { label: $t('system.tenant.status.disabled'), value: 0 },
      ],
    },
    defaultValue: 1,
    fieldName: 'status',
    label: $t('system.tenant.form.status'),
  },
];
