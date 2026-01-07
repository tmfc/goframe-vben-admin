import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'ant-design:setting-outlined',
      order: 10,
      title: $t('system.title'),
    },
    name: 'System',
    path: '/system',
    children: [
      {
        name: 'Menu',
        path: '/system/menu',
        component: () => import('#/views/sys/menu/index.vue'),
        meta: {
          icon: 'ant-design:menu-outlined',
          title: $t('system.menu.title'),
        },
      },
      {
        name: 'Dept',
        path: '/system/dept',
        component: () => import('#/views/sys/dept/index.vue'),
        meta: {
          icon: 'ant-design:cluster-outlined',
          title: $t('system.dept.title'),
        },
      },
    ],
  },
];

export default routes;
