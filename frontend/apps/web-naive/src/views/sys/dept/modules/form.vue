<script lang="ts" setup>
import type { Dept } from '#/api/sys/dept';

import { computed, ref, watch } from 'vue';

import { NButton, NModal, NSpace } from 'naive-ui';

import { $t } from '#/locales';

import { useVbenForm } from '#/adapter/form';
import { createDept, updateDept } from '#/api/sys/dept';

import { useFormSchema } from '../data';

const props = defineProps<{
  record: Dept | null;
  show: boolean;
}>();

const emit = defineEmits<{
  (event: 'success'): void;
  (event: 'update:show', value: boolean): void;
}>();

const saving = ref(false);

const [Form, formApi] = useVbenForm({
  schema: useFormSchema(),
  showDefaultActions: false,
});

const isUpdate = computed(() => Boolean(props.record?.id));
const modalTitle = computed(() =>
  isUpdate.value
    ? $t('system.dept.form.editTitle')
    : $t('system.dept.form.createTitle'),
);

watch(
  () => props.show,
  (show) => {
    if (!show) return;
    if (props.record) {
      formApi.setValues(props.record);
      return;
    }
    const resetForm = (formApi as any).resetForm;
    if (typeof resetForm === 'function') {
      resetForm();
      return;
    }
    formApi.setValues({
      name: '',
      parentId: null,
      order: 1,
      status: 1,
    });
  },
);

function closeModal() {
  emit('update:show', false);
}

async function handleSubmit() {
  const { valid } = await formApi.validate();
  if (!valid) return;

  const values = await formApi.getValues();
  const payload = {
    name: values.name,
    order: values.order,
    status: values.status,
    parentId: values.parentId ?? '',
  };
  saving.value = true;
  try {
    const recordId = props.record?.id ?? values.id;
    await (isUpdate.value
      ? updateDept(recordId, payload as any)
      : createDept(payload as any));
    emit('success');
    closeModal();
  } finally {
    saving.value = false;
  }
}
</script>
<template>
  <NModal
    :show="show"
    preset="dialog"
    :title="modalTitle"
    :mask-closable="false"
    @update:show="emit('update:show', $event)"
  >
    <Form />
    <template #action>
      <NSpace>
        <NButton @click="closeModal">{{ $t('common.cancel') }}</NButton>
        <NButton type="primary" :loading="saving" @click="handleSubmit">
          {{ $t('common.confirm') }}
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
