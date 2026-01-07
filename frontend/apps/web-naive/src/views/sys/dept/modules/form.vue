<script lang="ts" setup>
import { computed } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '#/locales';

import { useVbenForm } from '#/adapter/form';
import { createDept, updateDept } from '#/api/sys/dept';

import { useFormSchema } from '../data';

const emits = defineEmits(['success']);

const [Drawer, drawerApi] = useVbenDrawer();
const [Form, formApi] = useVbenForm({
  schema: useFormSchema(),
});

const getDrawerTitle = computed(() => {
  return drawerApi.isUpdate
    ? $t('system.dept.form.editTitle')
    : $t('system.dept.form.createTitle');
});

async function handleSubmit() {
  const { valid } = await formApi.validate();
  if (!valid) return;

  const values = await formApi.getValues();
  drawerApi.lock();
  try {
    await (drawerApi.isUpdate
      ? updateDept(values.id, values)
      : createDept(values));
    emits('success');
    drawerApi.close();
  } finally {
    drawerApi.unlock();
  }
}
</script>
<template>
  <Drawer :title="getDrawerTitle" @confirm="handleSubmit">
    <Form />
  </Drawer>
</template>
