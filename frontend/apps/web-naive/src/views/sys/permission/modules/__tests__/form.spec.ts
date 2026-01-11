import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import Form from '../form.vue';
import { defineComponent, h } from 'vue';

vi.mock('@vben/common-ui', () => ({}));

const formApi = {
  validate: vi.fn(async () => ({ valid: true })),
  getValues: vi.fn(async () => ({})),
  setValues: vi.fn(),
  resetForm: vi.fn(),
  setFieldValue: vi.fn(),
};

vi.mock('#/adapter/form', () => ({
  useVbenForm: vi.fn(() => {
    const FormStub = defineComponent({
      name: 'FormStub',
      setup() {
        return () => h('div');
      },
    });
    return [FormStub, formApi];
  }),
}));

vi.mock('#/api/sys/permission', () => ({
  createPermission: vi.fn(() => Promise.resolve()),
  updatePermission: vi.fn(() => Promise.resolve()),
}));

describe('Permission Form', () => {
  it('should mount successfully', () => {
    const wrapper = mount(Form, {
      props: {
        show: true,
        record: null,
      },
      global: {
        stubs: {
          teleport: true,
          NModal: defineComponent({
            props: { show: Boolean },
            setup(props, { slots }) {
              return () => props.show ? h('div', [slots.default?.(), slots.action?.()]) : null;
            }
          }),
          NSpace: defineComponent({
            setup(_props, { slots }) { return () => h('div', slots.default?.()); }
          }),
          NButton: defineComponent({
            setup(_props, { slots }) { return () => h('button', slots.default?.()); }
          }),
        }
      }
    });
    expect(wrapper.exists()).toBe(true);
  });
});
