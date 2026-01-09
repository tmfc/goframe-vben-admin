import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import RoleManagement from '../index.vue';
import { NConfigProvider, NDialogProvider } from 'naive-ui';

import { getRoleList } from '#/api/sys/role';

vi.mock('#/api/sys/role', () => ({
  getRoleList: vi.fn(() => Promise.resolve({ items: [], total: 0 })),
  deleteRole: vi.fn(() => Promise.resolve()),
  getRolePermissions: vi.fn(() => Promise.resolve({ permissionIds: [] })),
  assignRolePermissions: vi.fn(() => Promise.resolve()),
}));

vi.mock('#/api/sys/permission', () => ({
  getPermissionList: vi.fn(() => Promise.resolve({ items: [], total: 0 })),
}));

describe('RoleManagement', () => {
  const mountComponent = () => mount(
    {
      template: '<n-config-provider><n-dialog-provider><role-management /></n-dialog-provider></n-config-provider>',
      components: {
        RoleManagement,
        NConfigProvider,
        NDialogProvider,
      },
    },
    {
      global: {
        stubs: {
          teleport: true,
          VbenTable: {
            template: '<div></div>',
            methods: {
              reload: vi.fn(),
              setProps: vi.fn(),
            }
          },
          VbenButton: {
            template: '<button><slot></slot></button>'
          }
        },
      },
    },
  );

  it('should mount successfully', () => {
    const wrapper = mountComponent();
    expect(wrapper.exists()).toBe(true);
  });

  it('should request role list on mount', () => {
    mountComponent();
    expect(getRoleList).toHaveBeenCalled();
  });
});
