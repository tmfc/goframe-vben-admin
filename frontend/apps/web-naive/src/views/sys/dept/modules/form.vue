<script lang="ts" setup>
import { computed } from 'vue';

import { useDrawer } from '@vben/common-ui';

import { useVbenForm } from '#/adapter/form';
import { createDept, updateDept } from '#/api/sys/dept';

import { useFormSchema } from '../data';

const emits = defineEmits(['success']);

const [Drawer, drawerApi] = useDrawer();
const [Form, formApi] = useVbenForm({
  schema: useFormSchema(),
});

const getDrawerTitle = computed(() => {
  return drawerApi.isUpdate ? 'Edit Department' : 'Create Department';
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
