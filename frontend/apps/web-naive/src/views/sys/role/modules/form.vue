<script lang="ts" setup>
import type { Role } from '#/api/sys/role';

import { computed, ref, watch } from 'vue';

import { NButton, NModal, NSpace } from 'naive-ui';

import { useVbenForm } from '#/adapter/form';
import { createRole, updateRole } from '#/api/sys/role';
import { $t } from '#/locales';

import { useFormSchema } from '../data';

const props = defineProps<{
  record: Role | null;
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
    ? $t('system.role.form.editTitle')
    : $t('system.role.form.createTitle'),
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
      description: '',
      parentId: null,
      deptId: null,
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
    description: values.description,
    status: values.status,
    parent_id: values.parentId ?? 0,
    dept_id: values.deptId ?? 0,
  };
  saving.value = true;
  try {
    const recordId = props.record?.id ?? values.id;
    await (isUpdate.value
      ? updateRole(recordId, payload as any)
      : createRole(payload as any));
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
