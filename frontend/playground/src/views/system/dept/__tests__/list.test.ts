import { mount } from '@vue/test-utils';
import { Button } from 'ant-design-vue';
import { describe, expect, it, vi } from 'vitest';

import List from '../list.vue';

// Mock dependencies
vi.mock('#/adapter/vxe-table', () => ({
  useVbenVxeGrid: vi.fn(() => [
    {
      template: '<div><slot name="toolbar-tools"></slot></div>',
    },
    {
      query: vi.fn(),
    },
  ]),
}));

vi.mock('#/api/system/dept', () => ({
  getDeptList: vi.fn(() => Promise.resolve([])),
  deleteDept: vi.fn(() => Promise.resolve()),
}));

vi.mock('@vben/common-ui', async () => {
  const original = await vi.importActual('@vben/common-ui');
  return {
    ...original,
    useVbenModal: vi.fn(() => [
      { template: '<div>FormModal</div>' },
      {
        setData: vi.fn().mockReturnThis(),
        open: vi.fn(),
      },
    ]),
  };
});

vi.mock('#/locales', () => ({
  $t: (key: string) => key,
}));

// Mock ant-design-vue message
vi.mock('ant-design-vue', async () => {
  const original = await vi.importActual('ant-design-vue');
  return {
    ...original,
    message: {
      loading: vi.fn(() => vi.fn()),
      success: vi.fn(),
    },
  };
});

describe('DeptList', () => {
  it('should render the page title and create button', () => {
    const wrapper = mount(List);

    // Check for the title passed to the Grid component
    expect(wrapper.html()).toContain('table-title="部门列表"');

    // Check for the create button
    const createButton = wrapper.findComponent(Button);
    expect(createButton.exists()).toBe(true);
    const plusIcon = createButton.find('svg');
    expect(plusIcon.exists()).toBe(true);
  });
});
