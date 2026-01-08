import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import Form from '../form.vue';
import { NForm } from 'naive-ui';

vi.mock('@vben/common-ui', () => ({}));

vi.mock('#/adapter/form', () => ({
  useVbenForm: vi.fn(() => ([vi.fn(), {}])),
}));

describe('User Form', () => {
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
