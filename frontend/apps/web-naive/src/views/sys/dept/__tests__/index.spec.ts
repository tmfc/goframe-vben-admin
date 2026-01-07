import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import DeptManagement from '../index.vue';
import { NConfigProvider, NDialogProvider } from 'naive-ui';

vi.mock('#/api/sys/dept', () => ({
  createDept: vi.fn(),
  updateDept: vi.fn(),
  deleteDept: vi.fn(),
  getDeptList: vi.fn(() => Promise.resolve({ items: [], total: 0 })),
}));

vi.mock('@vben/common-ui', () => ({
  useVbenDrawer: vi.fn(() => ([{}, { open: vi.fn(), close: vi.fn(), isUpdate: false, lock: vi.fn(), unlock: vi.fn(), getData: vi.fn() }])),
}));

vi.mock('#/adapter/form', () => ({
  useVbenForm: vi.fn(() => ([vi.fn(), {}])),
}));

describe('DeptManagement', () => {
  const mountComponent = () => mount(
    {
      template: '<n-config-provider><n-dialog-provider><dept-management /></n-dialog-provider></n-config-provider>',
      components: {
        DeptManagement,
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
});
