<script lang="ts" setup>
import type { User } from '@vben/types';

import { computed, ref, watch } from 'vue';

import { NButton, NForm, NFormItem, NInput, NModal, NSpace } from 'naive-ui';

import { $t } from '#/locales';

import { useVbenForm } from '#/adapter/form';
import { createUser, updateUser } from '#/api/sys/user';

import { useFormSchema } from '../data';

const props = defineProps<{
  record: User | null;
  show: boolean;
}>();

const emit = defineEmits<{
  (event: 'success'): void;
  (event: 'update:show', value: boolean): void;
}>();

const saving = ref(false);
const showPasswordModal = ref(false);
const passwordValue = ref('');
const savingPassword = ref(false);

const [Form, formApi] = useVbenForm({
  schema: useFormSchema(() => isUpdate.value),
  showDefaultActions: false,
});

const isUpdate = computed(() => Boolean(props.record?.id));
const modalTitle = computed(() =>
  isUpdate.value ? $t('system.user.form.editTitle') : $t('system.user.form.createTitle'),
);

function normalizeRolesToInput(roles: unknown): string[] {
  if (!roles) return [];
  if (Array.isArray(roles)) {
    return roles.map((item) => String(item)).filter(Boolean);
  }
  if (typeof roles === 'string') {
    const trimmed = roles.trim();
    if (!trimmed) return [];
    if (trimmed.startsWith('[')) {
      try {
        const parsed = JSON.parse(trimmed);
        if (Array.isArray(parsed)) {
          return parsed.map((item) => String(item)).filter(Boolean);
        }
      } catch {
        // fall through to comma parsing
      }
    }
    return trimmed
      .split(',')
      .map((item) => item.trim())
      .filter(Boolean);
  }
  return [];
}

function normalizeRolesForSubmit(input: unknown) {
  if (!input) return '';
  if (Array.isArray(input)) {
    const list = input.map((item) => String(item)).filter(Boolean);
    return list.length ? JSON.stringify(list) : '';
  }
  if (typeof input === 'string') {
    const trimmed = input.trim();
    if (!trimmed) return '';
    if (trimmed.startsWith('[')) {
      try {
        const parsed = JSON.parse(trimmed);
        if (Array.isArray(parsed)) {
          const list = parsed.map((item) => String(item)).filter(Boolean);
          return list.length ? JSON.stringify(list) : '';
        }
      } catch {
        // fall through to comma parsing
      }
    }
    const list = trimmed
      .split(',')
      .map((item) => item.trim())
      .filter(Boolean);
    return list.length ? JSON.stringify(list) : '';
  }
  return '';
}

watch(
  () => props.show,
  (show) => {
    if (!show) return;
    const updateSchema = (formApi as any).updateSchema;
    if (typeof updateSchema === 'function') {
      updateSchema([
        {
          fieldName: 'password',
          hide: isUpdate.value,
          rules: isUpdate.value ? null : 'required',
        },
      ]);
    }
    if (props.record) {
      formApi.setValues({
        ...props.record,
        id: props.record.id,
        roles: normalizeRolesToInput((props.record as any).roles),
      });
      return;
    }
    const resetForm = (formApi as any).resetForm;
    if (typeof resetForm === 'function') {
      resetForm();
      return;
    }
      formApi.setValues({
        id: '',
        username: '',
        realName: '',
        password: '',
        roles: [],
        homePath: '',
        avatar: '',
        deptId: null,
      status: 1,
    });
  },
);

function closeModal() {
  emit('update:show', false);
}

function openPasswordModal() {
  passwordValue.value = '';
  showPasswordModal.value = true;
}

async function submitPasswordChange() {
  if (!props.record?.id || !passwordValue.value) {
    return;
  }
  savingPassword.value = true;
  try {
    await updateUser(props.record.id, {
      username: props.record.username,
      status: props.record.status,
      password: passwordValue.value,
    });
    showPasswordModal.value = false;
  } finally {
    savingPassword.value = false;
  }
}

async function handleSubmit() {
  const { valid } = await formApi.validate();
  if (!valid) return;

  const values = await formApi.getValues();
  saving.value = true;
  try {
    const recordId = props.record?.id ?? values.id;
    const payload = {
      ...values,
      roles: normalizeRolesForSubmit(values.roles),
    };
    await (isUpdate.value
      ? updateUser(recordId, payload)
      : createUser(payload));
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
        <NButton v-if="isUpdate" type="warning" @click="openPasswordModal">
          {{ $t('system.user.actions.changePassword') }}
        </NButton>
        <NButton @click="closeModal">{{ $t('common.cancel') }}</NButton>
        <NButton type="primary" :loading="saving" @click="handleSubmit">
          {{ $t('common.confirm') }}
        </NButton>
      </NSpace>
    </template>
  </NModal>
  <NModal
    v-model:show="showPasswordModal"
    preset="dialog"
    :title="$t('system.user.password.title')"
    :mask-closable="false"
  >
    <NForm>
      <NFormItem :label="$t('system.user.password.newPassword')">
        <NInput
          v-model:value="passwordValue"
          type="password"
          show-password-on="click"
        />
      </NFormItem>
    </NForm>
    <template #action>
      <NSpace>
        <NButton @click="showPasswordModal = false">
          {{ $t('common.cancel') }}
        </NButton>
        <NButton
          type="primary"
          :loading="savingPassword"
          :disabled="!passwordValue"
          @click="submitPasswordChange"
        >
          {{ $t('common.confirm') }}
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
