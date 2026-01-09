<script setup lang="ts">
import { computed, ref, watch } from 'vue';

import { usePagination } from '@vben/hooks';
import { EmptyIcon, Grip, IconifyIcon, listIcons } from '@vben/icons';
import { refDebounced } from '@vueuse/core';

import { NButton, NInput, NPagination, NPopover } from 'naive-ui';

import { $t } from '#/locales';

interface Props {
  modelValue?: string;
  prefix?: string;
  pageSize?: number;
  placeholder?: string;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  prefix: '',
  pageSize: 36,
  placeholder: '',
});

const emit = defineEmits<{
  'update:modelValue': [string];
}>();

const show = ref(false);
const keyword = ref('');
const keywordDebounced = refDebounced(keyword, 300);
const current = ref(props.modelValue);
const activePrefix = ref('');
const innerIcons = ref<string[]>([]);

const iconsCache: Record<string, string[]> = {};
const pendingRequests: Record<string, Promise<string[]>> = {};

async function fetchIconsData(prefix: string): Promise<string[]> {
  if (iconsCache[prefix]?.length) {
    return iconsCache[prefix];
  }
  if (pendingRequests[prefix]) {
    return pendingRequests[prefix];
  }
  pendingRequests[prefix] = (async () => {
    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 1000 * 10);
      const response = await fetch(
        `https://api.iconify.design/collection?prefix=${prefix}`,
        { signal: controller.signal },
      ).then((res) => res.json());
      clearTimeout(timeoutId);
      const list = response?.uncategorized || [];
      if (response?.categories) {
        for (const category in response.categories) {
          list.push(...(response.categories[category] || []));
        }
      }
      iconsCache[prefix] = list.map((value: string) => `${prefix}:${value}`);
    } catch (error) {
      console.error(`Failed to fetch icons for prefix ${prefix}:`, error);
      iconsCache[prefix] = [];
    }
    return iconsCache[prefix];
  })();
  return pendingRequests[prefix];
}

function resolvePrefix(value: string) {
  if (props.prefix) {
    return props.prefix;
  }
  if (value?.includes(':')) {
    return value.split(':')[0];
  }
  return 'ant-design';
}

function resolvePrefixFromKeyword(value: string) {
  const match = value.match(/^([a-z0-9-]+):/i);
  return match?.[1] ?? '';
}

watch(
  () => props.modelValue,
  (value) => {
    current.value = value ?? '';
    activePrefix.value = resolvePrefix(current.value);
  },
  { immediate: true },
);

watch(
  () => props.prefix,
  () => {
    activePrefix.value = resolvePrefix(current.value);
  },
);

watch(
  () => keywordDebounced.value,
  (value) => {
    if (props.prefix) {
      return;
    }
    const nextPrefix = resolvePrefixFromKeyword(value);
    if (nextPrefix && nextPrefix !== activePrefix.value) {
      activePrefix.value = nextPrefix;
    }
  },
);

watch(current, (value) => {
  if (value !== props.modelValue) {
    emit('update:modelValue', value);
  }
});

watch(
  () => activePrefix.value,
  async (prefix) => {
    if (!prefix) {
      innerIcons.value = [];
      return;
    }
    innerIcons.value = await fetchIconsData(prefix);
  },
  { immediate: true },
);

const allIcons = computed(() => {
  const prefix = activePrefix.value?.trim();
  if (!prefix) {
    return [];
  }
  const icons = listIcons('', prefix);
  if (innerIcons.value.length > 0) {
    return innerIcons.value;
  }
  return icons;
});

const filteredIcons = computed(() => {
  const keywordValue = keywordDebounced.value.trim();
  if (!keywordValue) {
    return allIcons.value;
  }
  return allIcons.value.filter((item) => item.includes(keywordValue));
});

const { paginationList, total, setCurrentPage, currentPage } = usePagination(
  filteredIcons,
  props.pageSize,
);

const totalPages = computed(() => {
  return Math.ceil(total.value / props.pageSize);
});

const triggerPlaceholder = computed(() => {
  return props.placeholder || $t('ui.iconPicker.placeholder');
});

function handleSelect(icon: string) {
  current.value = icon;
  show.value = false;
}

function handlePageChange(page: number) {
  setCurrentPage(page);
}
</script>

<template>
  <NPopover
    :show="show"
    trigger="manual"
    placement="bottom-end"
    :width="360"
    :show-arrow="false"
    :to="false"
    @clickoutside="show = false"
  >
    <template #trigger>
      <NInput
        v-model:value="current"
        :placeholder="triggerPlaceholder"
        readonly
        @click="show = true"
      >
        <template #suffix>
          <IconifyIcon v-if="current" :icon="current" class="size-4" />
          <Grip v-else class="size-4" />
        </template>
      </NInput>
    </template>

    <div class="menu-icon-picker" @mousedown.stop @click.stop>
      <NInput
        v-model:value="keyword"
        class="menu-icon-picker__search"
        :placeholder="$t('ui.iconPicker.search')"
        clearable
      />

      <template v-if="paginationList.length > 0">
        <div class="menu-icon-picker__grid">
          <NButton
            v-for="item in paginationList"
            :key="item"
            size="small"
            quaternary
            class="menu-icon-picker__item"
            @click="handleSelect(item)"
          >
            <IconifyIcon
              :icon="item"
              :class="{
                'menu-icon-picker__active': current === item,
              }"
            />
          </NButton>
        </div>
        <div v-if="totalPages > 1" class="menu-icon-picker__pager">
          <NPagination
            :page="currentPage"
            :page-count="totalPages"
            :page-slot="5"
            size="small"
            @update:page="handlePageChange"
          />
        </div>
      </template>

      <div v-else class="menu-icon-picker__empty">
        <EmptyIcon class="menu-icon-picker__empty-icon" />
        <div class="menu-icon-picker__empty-text">{{ $t('common.noData') }}</div>
      </div>
    </div>
  </NPopover>
</template>

<style scoped>
.menu-icon-picker {
  padding: 8px;
}

.menu-icon-picker__search {
  margin-bottom: 8px;
}

.menu-icon-picker__grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 6px;
  max-height: 360px;
  overflow: auto;
  padding: 4px;
}

.menu-icon-picker__item {
  display: flex;
  justify-content: center;
}

.menu-icon-picker__active {
  color: var(--n-primary-color);
}

.menu-icon-picker__pager {
  display: flex;
  justify-content: flex-end;
  padding-top: 8px;
  border-top: 1px solid var(--n-border-color);
  overflow-x: auto;
  white-space: nowrap;
}

:deep(.n-pagination) {
  flex-wrap: nowrap;
  white-space: nowrap;
}

:deep(.n-pagination-item) {
  white-space: nowrap;
}

.menu-icon-picker__empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 0;
  color: var(--n-text-color-3);
}

.menu-icon-picker__empty-icon {
  width: 32px;
  height: 32px;
}

.menu-icon-picker__empty-text {
  margin-top: 4px;
  font-size: 12px;
}
</style>
