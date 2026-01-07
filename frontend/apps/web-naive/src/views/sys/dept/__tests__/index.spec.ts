import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import DeptManagement from '../index.vue';
import { createDept, updateDept, deleteDept } from '#/api/sys/dept';

vi.mock('#/api/sys/dept', () => ({
  createDept: vi.fn(),
  updateDept: vi.fn(),
  deleteDept: vi.fn(),
}));

describe('DeptManagement', () => {
  it('should mount successfully', () => {
    const wrapper = mount(DeptManagement);
    expect(wrapper.exists()).toBe(true);
  });

  it('should create a new department', async () => {
    const wrapper = mount(DeptManagement);
    await wrapper.find('button.create-btn').trigger('click');
    await wrapper.find('input.name-input').setValue('Test Dept');
    await wrapper.find('button.submit-btn').trigger('click');
    expect(createDept).toHaveBeenCalledWith({ name: 'Test Dept' });
  });

  it('should edit a department', async () => {
    const wrapper = mount(DeptManagement);
    await wrapper.find('button.edit-btn').trigger('click');
    await wrapper.find('input.name-input').setValue('Updated Dept');
    await wrapper.find('button.submit-btn').trigger('click');
    expect(updateDept).toHaveBeenCalledWith(expect.any(String), { name: 'Updated Dept' });
  });

  it('should delete a department', async () => {
    const wrapper = mount(DeptManagement);
    await wrapper.find('button.delete-btn').trigger('click');
    expect(deleteDept).toHaveBeenCalledWith(expect.any(String));
  });
});