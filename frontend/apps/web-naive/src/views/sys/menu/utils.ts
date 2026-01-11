import type { MenuItem } from '#/api/sys/menu';
import { $t } from '#/locales';
import { ref, watch } from 'vue';
import { useI18n } from '@vben/locales';
import zhSystem from '#/locales/langs/zh-CN/system.json';
import enSystem from '#/locales/langs/en-US/system.json';

export function getMeta(row: any) {
  if (!row) return null;
  const rawMeta = row.meta;
  if (typeof rawMeta === 'string') {
    try {
      return JSON.parse(rawMeta);
    } catch {
      return null;
    }
  }
  if (rawMeta && typeof rawMeta === 'object') return rawMeta;
  return null;
}

export function getMenuTitle(row: any) {
  if (!row) return '';
  const meta = getMeta(row);
  const titleKey = meta?.title;
  if (titleKey) {
    return $t(titleKey);
  }
  return row.name ?? '';
}

export function useTitleOptions() {
  const { locale } = useI18n();
  const options = ref<{ label: string; value: string }[]>([]);
  
  const localeSystemMap: Record<string, any> = {
    'zh-CN': zhSystem,
    'en-US': enSystem,
  };

  function rebuild() {
    const set = new Set<string>();
    function collect(obj: any, prefix = '') {
      if (!obj || typeof obj !== 'object') return;
      Object.entries(obj).forEach(([k, v]) => {
        const keyPath = prefix ? `${prefix}${k}` : k;
        if (typeof v === 'string') {
          set.add(keyPath);
          return;
        }
        if (v && typeof v === 'object') {
          collect(v, `${keyPath}.`);
        }
      });
    }
    const currentLocale = locale.value;
    const currentSystem = localeSystemMap[currentLocale] || localeSystemMap['zh-CN'];
    collect(currentSystem?.common, 'common.');
    collect(currentSystem?.menu, 'system.menu.');
    collect(currentSystem?.permission, 'system.permission.');
    options.value = Array.from(set).map((key) => ({
      label: `${$t(key)} (${key})`,
      value: key,
    }));
  }

  watch(() => locale.value, rebuild, { immediate: true });

  return options;
}

export function generatePermissionCode(
  name: string,
  type: string,
  parentId?: string | null,
  rawMenuList: MenuItem[] = [],
): string {
  // 将菜单名称转换为帕斯卡命名 (首字母大写)
  const pascalCaseName = name
    .replace(/[-_\s]+(.)?/g, (_: string, char: string) =>
      char ? char.toUpperCase() : '',
    )
    .replace(/^(.)/, (_: string, char: string) => char.toUpperCase());

  // 根据类型生成不同的前缀
  const prefixMap: Record<string, string> = {
    button: 'Button',
    catalog: 'Catalog',
    embedded: 'Embedded',
    link: 'Link',
    menu: 'Menu',
  };

  const prefix = prefixMap[type] || 'Menu';

  // 如果是按钮类型,需要获取父菜单的路径
  if (type === 'button' && parentId) {
    const parentMenu = rawMenuList.find(
      (item) => String(item.id) === String(parentId),
    );
    if (parentMenu && parentMenu.permissionCode) {
      // 从父菜单的权限代码中提取路径部分,去掉最后的操作词
      const parentPath = parentMenu.permissionCode;
      return `${parentPath}:${pascalCaseName}`;
    }
  }

  // 如果是菜单类型,需要获取完整的父级路径
  if (type !== 'button' && parentId) {
    const pathParts: string[] = [];
    let currentId = String(parentId);

    // 向上遍历父级菜单
    while (currentId && currentId !== '0') {
      const parent = rawMenuList.find((item) => String(item.id) === currentId);
      if (!parent) break;

      const pascalParentName = (parent.name || '')
        .replace(/[-_\s]+(.)?/g, (_: string, char: string) =>
          char ? char.toUpperCase() : '',
        )
        .replace(/^(.)/, (_: string, char: string) => char.toUpperCase());

      pathParts.unshift(pascalParentName);
      currentId = parent.parentId ? String(parent.parentId) : '';
    }

    // 构建完整路径: Menu:System:Permission:List
    const fullPath =
      pathParts.length > 0
        ? `${prefix}:${pathParts.join(':')}:${pascalCaseName}`
        : `${prefix}:${pascalCaseName}`;
    return fullPath;
  }

  return `${prefix}:${pascalCaseName}`;
}
