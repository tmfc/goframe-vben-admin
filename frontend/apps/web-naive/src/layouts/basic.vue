<script lang="ts" setup>
import { computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';

import { AuthenticationLoginExpiredModal } from '@vben/common-ui';
import { VBEN_DOC_URL, VBEN_GITHUB_URL } from '@vben/constants';
import { useWatermark } from '@vben/hooks';
import { BookOpenText, CircleHelp, SvgGithubIcon } from '@vben/icons';
import {
  BasicLayout,
  LockScreen,
  Notification,
  UserDropdown,
} from '@vben/layouts';
import { preferences } from '@vben/preferences';
import { useAccessStore, useUserStore } from '@vben/stores';
import { openWindow } from '@vben/utils';
import { NButton, NForm, NFormItem, NModal, NSelect, NSpace } from 'naive-ui';

import { $t } from '#/locales';
import { getTenantList } from '#/api/sys/tenant';
import { useAuthStore, useMessageStore } from '#/store';
import LoginForm from '#/views/_core/authentication/login.vue';

import { useNotification } from './useNotification';

const router = useRouter();
const userStore = useUserStore();
const authStore = useAuthStore();
const accessStore = useAccessStore();
const messageStore = useMessageStore();
const { destroyWatermark, updateWatermark } = useWatermark();

const {
  handleMakeAll,
  handleNoticeClear,
  markRead,
  notifications,
  remove,
  showDot,
} = useNotification();
const multiTenantEnabled = import.meta.env.VITE_APP_MULTI_TENANT === 'true';

const menus = computed(() => {
  const items = [
    {
      handler: () => {
        router.push({ name: 'Profile' });
      },
      icon: 'lucide:user',
      text: $t('page.auth.profile'),
    },
    {
      handler: () => {
        openWindow(VBEN_DOC_URL, {
          target: '_blank',
        });
      },
      icon: BookOpenText,
      text: $t('ui.widgets.document'),
    },
    {
      handler: () => {
        openWindow(VBEN_GITHUB_URL, {
          target: '_blank',
        });
      },
      icon: SvgGithubIcon,
      text: 'GitHub',
    },
    {
      handler: () => {
        openWindow(`${VBEN_GITHUB_URL}/issues`, {
          target: '_blank',
        });
      },
      icon: CircleHelp,
      text: $t('ui.widgets.qa'),
    },
  ];

  if (multiTenantEnabled && userStore.userRoles.includes('super')) {
    items.unshift({
      handler: () => {
        openTenantSwitch();
      },
      icon: 'lucide:repeat',
      text: 'Switch Tenant',
    });
  }

  return items;
});

const avatar = computed(() => {
  return userStore.userInfo?.avatar ?? preferences.app.defaultAvatar;
});

async function handleLogout() {
  await authStore.logout(false);
}

onMounted(() => {
  messageStore.connect();
});

const tenantSwitchVisible = ref(false);
const tenantLoading = ref(false);
const switchingTenant = ref(false);
const selectedTenantId = ref<number | string>();
const tenantOptions = ref<{ label: string; value: number | string }[]>([]);

async function openTenantSwitch() {
  tenantSwitchVisible.value = true;
  selectedTenantId.value = undefined;
  if (tenantOptions.value.length > 0) {
    return;
  }
  tenantLoading.value = true;
  try {
    const res = await getTenantList({ page: 1, pageSize: 200 });
    const items = res?.items ?? res?.data?.items ?? [];
    tenantOptions.value = items.map((item: any) => ({
      label: `${item.name} (${item.code || item.id})`,
      value: item.id,
    }));
  } finally {
    tenantLoading.value = false;
  }
}

async function handleTenantSwitch() {
  if (!selectedTenantId.value) {
    return;
  }
  switchingTenant.value = true;
  try {
    await authStore.switchTenant(selectedTenantId.value);
    tenantSwitchVisible.value = false;
  } finally {
    switchingTenant.value = false;
  }
}

watch(
  () => ({
    enable: preferences.app.watermark,
    content: preferences.app.watermarkContent,
  }),
  async ({ enable, content }) => {
    if (enable) {
      await updateWatermark({
        content:
          content ||
          `${userStore.userInfo?.username} - ${userStore.userInfo?.realName}`,
      });
    } else {
      destroyWatermark();
    }
  },
  {
    immediate: true,
  },
);
</script>

<template>
  <BasicLayout @clear-preferences-and-logout="handleLogout">
    <template #user-dropdown>
      <UserDropdown
        :avatar
        :menus
        :text="userStore.userInfo?.realName"
        description="ann.vben@gmail.com"
        tag-text="Pro"
        @logout="handleLogout"
      />
    </template>
    <template #notification>
      <Notification
        :dot="showDot"
        :notifications="notifications"
        @clear="handleNoticeClear"
        @read="(item) => item.id && markRead(item.id)"
        @remove="(item) => item.id && remove(item.id)"
        @make-all="handleMakeAll"
      />
    </template>
    <template #extra>
      <AuthenticationLoginExpiredModal
        v-model:open="accessStore.loginExpired"
        :avatar
      >
        <LoginForm />
      </AuthenticationLoginExpiredModal>
    </template>
    <template #lock-screen>
      <LockScreen :avatar @to-login="handleLogout" />
    </template>
    <NModal
      v-model:show="tenantSwitchVisible"
      preset="card"
      title="Switch Tenant"
      style="width: 420px"
    >
      <NForm>
        <NFormItem label="Tenant">
          <NSelect
            v-model:value="selectedTenantId"
            :loading="tenantLoading"
            :options="tenantOptions"
            filterable
            placeholder="Select tenant"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="tenantSwitchVisible = false">Cancel</NButton>
          <NButton
            type="primary"
            :disabled="!selectedTenantId"
            :loading="switchingTenant"
            @click="handleTenantSwitch"
          >
            Confirm
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </BasicLayout>
</template>
