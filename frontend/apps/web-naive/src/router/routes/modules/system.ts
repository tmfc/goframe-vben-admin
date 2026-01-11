import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'ant-design:setting-outlined',
      order: 10,
      title: 'system.title',
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
          title: 'system.menu.title',
        },
      },
      {
        name: 'Dept',
        path: '/system/dept',
        component: () => import('#/views/sys/dept/index.vue'),
        meta: {
          icon: 'ant-design:cluster-outlined',
          title: 'system.dept.title',
        },
      },
      {
        name: 'Permission',
        path: '/system/permission',
        component: () => import('#/views/sys/permission/index.vue'),
        meta: {
          icon: 'ant-design:safety-outlined',
          title: 'system.permission.title',
        },
      },
      {
        name: 'User',
        path: '/system/user',
        component: () => import('#/views/sys/user/index.vue'),
        meta: {
          icon: 'ant-design:user-outlined',
          title: 'system.user.title',
        },
      },
      {
        name: 'Tenant',
        path: '/system/tenant',
        component: () => import('#/views/sys/tenant/index.vue'),
        meta: {
          icon: 'ant-design:solution-outlined',
          title: 'system.tenant.title',
        },
      },
    ],
  },
];

export default routes;
