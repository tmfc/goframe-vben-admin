<script lang="ts" setup>
import { computed, ref, watch } from 'vue';

import { NButton, NModal, NSpace } from 'naive-ui';

import { useVbenForm } from '#/adapter/form';
import {
  createPermission,
  updatePermission,
  type PermissionItem,
} from '#/api/sys/permission';
import { $t } from '#/locales';

import { useFormSchema } from '../data';

const props = defineProps<{
  record: PermissionItem | null;
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
    ? $t('system.permission.actions.edit')
    : $t('system.permission.actions.create'),
);

watch(
  () => props.show,
  (show) => {
    if (!show) return;
    if (props.record) {
      formApi.setValues({
        ...props.record,
        parentId: props.record.parentId ? String(props.record.parentId) : '0',
      });
    } else {
      formApi.setValues({
        description: '',
        id: '',
        name: '',
        parentId: '0',
        status: 1,
      });
    }
  },
);

function closeModal() {
  emit('update:show', false);
}

async function handleSubmit() {
  const { valid } = await formApi.validate();
  if (!valid) return;

  const values = await formApi.getValues();
  saving.value = true;
  try {
    const parentIdValue = values.parentId ? Number(values.parentId) : 0;
    const payload = {
      ...values,
      parent_id: Number.isFinite(parentIdValue) ? parentIdValue : 0,
    } as any;
    await (isUpdate.value
      ? updatePermission(values.id, payload)
      : createPermission(payload));
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
    style="width: 600px"
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
