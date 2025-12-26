import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'ant-design:setting-outlined',
      order: 10,
      title: 'System',
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
          title: 'Menu Management',
        },
      },
    ],
  },
];

export default routes;
