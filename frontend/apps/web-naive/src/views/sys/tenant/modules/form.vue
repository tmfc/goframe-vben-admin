<script lang="ts" setup>
import type { TenantModel } from '#/api/sys/tenant';

import { computed, ref, watch } from 'vue';

import { NButton, NModal, NSpace } from 'naive-ui';

import { $t } from '#/locales';

import { useVbenForm } from '#/adapter/form';
import { createTenant, updateTenant } from '#/api/sys/tenant';

import { useFormSchema } from '../data';

const props = defineProps<{
  record: TenantModel | null;
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
  isUpdate.value ? $t('system.tenant.form.editTitle') : $t('system.tenant.form.createTitle'),
);

watch(
  () => props.show,
  (show) => {
    if (!show) return;
    if (props.record) {
      formApi.setValues({
        ...props.record,
      });
      return;
    }
    const resetForm = (formApi as any).resetForm;
    if (typeof resetForm === 'function') {
      resetForm();
    } else {
      formApi.setValues({
        name: '',
        code: '',
        contactName: '',
        contactMobile: '',
        startAt: null,
        expireAt: null,
        packageVersion: '',
        domain: '',
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
    const recordId = props.record?.id;
    if (isUpdate.value && recordId) {
      await updateTenant(recordId, values as any);
    } else {
      await createTenant(values as any);
    }
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
