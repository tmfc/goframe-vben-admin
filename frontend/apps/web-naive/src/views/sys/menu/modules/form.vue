<script lang="ts" setup>
import { computed, ref, watch } from 'vue';

import { NButton, NModal, NSpace } from 'naive-ui';

import { useVbenForm } from '#/adapter/form';
import { createMenu, getMenuList, updateMenu } from '#/api/sys/menu';
import type { MenuItem } from '#/api/sys/menu';
import { $t } from '#/locales';

import { useFormSchema } from '../data';
import { generatePermissionCode } from '../utils';

const props = defineProps<{
  record: MenuItem | null;
  show: boolean;
}>();

const emit = defineEmits<{
  (event: 'success'): void;
  (event: 'update:show', value: boolean): void;
}>();

const saving = ref(false);
const rawMenuList = ref<MenuItem[]>([]);

const [Form, formApi] = useVbenForm({
  handleValuesChange: (values) => {
    const { name, type, autoGeneratePermission, parentId, permissionCode } = values;
    if (autoGeneratePermission && name && type) {
      const generated = generatePermissionCode(
        String(name),
        String(type),
        parentId as string,
        rawMenuList.value,
      );
      if (generated !== permissionCode) {
        formApi.setFieldValue('permissionCode', generated);
      }
    }
  },
  schema: useFormSchema(),
  showDefaultActions: false,
});

const isUpdate = computed(() => Boolean(props.record?.id));
const modalTitle = computed(() =>
  isUpdate.value
    ? $t('system.menu.actions.edit')
    : $t('system.menu.actions.create'),
);

async function fetchRawMenuList() {
  const res = await getMenuList();
  rawMenuList.value = res.list || [];
}

watch(
  () => props.show,
  async (show) => {
    if (!show) return;
    await fetchRawMenuList();
    if (props.record) {
      formApi.setValues({
        ...props.record,
        autoGeneratePermission: !props.record.permissionCode,
      });
    } else {
      formApi.setValues({
        autoGeneratePermission: true,
        component: '',
        icon: '',
        id: '',
        name: '',
        order: 0,
        parentId: '0',
        path: '',
        permissionCode: '',
        status: 1,
        type: 'menu',
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
    const data = {
      ...values,
      parentId: values.parentId || '0',
    } as any;
    await (isUpdate.value
      ? updateMenu(values.id, data)
      : createMenu(data));
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
    style="width: 800px"
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
