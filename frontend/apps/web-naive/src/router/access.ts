import type {
  ComponentRecordType,
  GenerateMenuAndRoutesOptions,
} from '@vben/types';

import { generateAccessible } from '@vben/access';
import { preferences } from '@vben/preferences';
import type { RouteRecordStringComponent } from '@vben/types';

import { message } from '#/adapter/naive';
import { getAllMenusApi } from '#/api';
import { BasicLayout, IFrameView } from '#/layouts';
import { $t } from '#/locales';

const forbiddenComponent = () => import('#/views/_core/fallback/forbidden.vue');

function normalizeMenuMeta(
  menus: RouteRecordStringComponent[],
): RouteRecordStringComponent[] {
  return (menus || []).map((menu) => {
    const rawIcon = (menu as any).icon as string | undefined;
    const meta = { ...(menu.meta ?? {}) };
    if (rawIcon && !meta.icon) {
      meta.icon = rawIcon;
    }
    return {
      ...menu,
      meta,
      children: menu.children ? normalizeMenuMeta(menu.children) : undefined,
    };
  });
}

async function generateAccess(options: GenerateMenuAndRoutesOptions) {
  const pageMap: ComponentRecordType = import.meta.glob('../views/**/*.vue');

  const layoutMap: ComponentRecordType = {
    BasicLayout,
    IFrameView,
  };

  return await generateAccessible(preferences.app.accessMode, {
    ...options,
    fetchMenuListAsync: async () => {
      message.loading(`${$t('common.loadingMenu')}...`, {
        duration: 1.5,
      });
      const menus = await getAllMenusApi();
      return normalizeMenuMeta(menus);
    },
    // 可以指定没有权限跳转403页面
    forbiddenComponent,
    // 如果 route.meta.menuVisibleWithForbidden = true
    layoutMap,
    pageMap,
  });
}

export { generateAccess };
