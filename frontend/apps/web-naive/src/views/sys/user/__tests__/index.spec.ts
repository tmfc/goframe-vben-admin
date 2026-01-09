import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import UserManagement from '../index.vue';
import { NConfigProvider, NDialogProvider } from 'naive-ui';

import { getUserList } from '#/api/sys/user';

vi.mock('#/api/sys/user', () => ({
  getUserList: vi.fn(() => Promise.resolve({ items: [], total: 0 })),
  deleteUser: vi.fn(() => Promise.resolve()),
}));

describe('UserManagement', () => {
  const mountComponent = () => mount(
    {
      template: '<n-config-provider><n-dialog-provider><user-management /></n-dialog-provider></n-config-provider>',
      components: {
        UserManagement,
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

  it('should request user list on mount', () => {
    mountComponent();
    expect(getUserList).toHaveBeenCalled();
  });
});
