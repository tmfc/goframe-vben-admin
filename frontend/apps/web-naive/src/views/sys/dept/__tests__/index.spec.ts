import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import DeptManagement from '../index.vue';
import { createDept, updateDept, deleteDept } from '#/api/sys/dept';
import { NConfigProvider, NDialogProvider } from 'naive-ui';

vi.mock('#/api/sys/dept', () => ({
  createDept: vi.fn(),
  updateDept: vi.fn(),
  deleteDept: vi.fn(),
}));

vi.mock('@vben/common-ui', () => ({
  useDrawer: vi.fn(() => ([{}, { open: vi.fn(), close: vi.fn(), isUpdate: false, lock: vi.fn(), unlock: vi.fn(), getData: vi.fn() }])),
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
        },
      },
    },
  );

  it('should mount successfully', () => {
    const wrapper = mountComponent();
    expect(wrapper.exists()).toBe(true);
  });

  it('should create a new department', async () => {
    const wrapper = mountComponent();
    await wrapper.find('button.create-btn').trigger('click');
    // Assuming you have an input with class 'name-input' inside the form
    await wrapper.find('input.name-input').setValue('Test Dept');
    await wrapper.find('button.submit-btn').trigger('click');
    expect(createDept).toHaveBeenCalledWith({ name: 'Test Dept' });
  });

  it('should edit a department', async () => {
    const wrapper = mountComponent();
    await wrapper.find('button.edit-btn').trigger('click');
    // Assuming you have an input with class 'name-input' inside the form
    await wrapper.find('input.name.input').setValue('Updated Dept');
    await wrapper.find('button.submit-btn').trigger('click');
    expect(updateDept).toHaveBeenCalledWith(expect.any(String), { name: 'Updated Dept' });
  });

  it('should delete a department', async () => {
    const wrapper = mountComponent();
    await wrapper.find('button.delete-btn').trigger('click');
    expect(deleteDept).toHaveBeenCalledWith(expect.any(String));
  });
});