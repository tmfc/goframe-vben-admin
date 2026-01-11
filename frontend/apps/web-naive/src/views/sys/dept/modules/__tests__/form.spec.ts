import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import Form from '../form.vue';
import { NForm } from 'naive-ui';

vi.mock('@vben/common-ui', () => ({
  useVbenDrawer: vi.fn(() => ([{}, { open: vi.fn(), close: vi.fn(), isUpdate: false, lock: vi.fn(), unlock: vi.fn(), getData: vi.fn() }])),
}));

vi.mock('#/adapter/form', () => ({
  useVbenForm: vi.fn(() => ([vi.fn(), {}])),
  z: {
    number: () => ({
      int: () => ({
        min: () => ({}),
      }),
    }),
  },
}));

describe('Department Form', () => {
  const mountComponent = () =>
    mount(Form, {
      props: {
        show: true,
        record: null,
      },
      global: {
        components: {
          NForm,
        },
      },
    });

  it('should mount successfully', () => {
    const wrapper = mountComponent();
    expect(wrapper.exists()).toBe(true);
  });
});
