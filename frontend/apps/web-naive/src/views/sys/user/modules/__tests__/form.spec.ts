import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import Form from '../form.vue';
import { NForm } from 'naive-ui';
import { defineComponent, h, nextTick } from 'vue';

import { createUser, updateUser } from '#/api/sys/user';

vi.mock('@vben/common-ui', () => ({}));

const formApi = {
  validate: vi.fn(async () => ({ valid: true })),
  getValues: vi.fn(async () => ({})),
  setValues: vi.fn(),
  resetForm: vi.fn(),
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

vi.mock('#/api/sys/user', () => ({
  createUser: vi.fn(() => Promise.resolve()),
  updateUser: vi.fn(() => Promise.resolve()),
}));

describe('User Form', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

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
        stubs: {
          NButton: defineComponent({
            props: {
              onClick: Function,
            },
            setup(props, { slots }) {
              return () =>
                h(
                  'button',
                  {
                    type: 'button',
                    onClick: props.onClick as any,
                  },
                  slots.default?.(),
                );
            },
          }),
          NModal: defineComponent({
            props: {
              show: Boolean,
            },
            setup(props, { slots }) {
              return () =>
                props.show ? h('div', [slots.default?.(), slots.action?.()]) : null;
            },
          }),
          NFormItem: defineComponent({
            setup(_props, { slots }) {
              return () => h('div', slots.default?.());
            },
          }),
          NInput: defineComponent({
            props: {
              value: String,
              'onUpdate:value': Function,
            },
            setup(props) {
              return () =>
                h('input', {
                  value: props.value,
                  onInput: (event: Event) => {
                    const target = event.target as HTMLInputElement;
                    (props['onUpdate:value'] as any)?.(target.value);
                  },
                });
            },
          }),
          NSpace: defineComponent({
            setup(_props, { slots }) {
              return () => h('div', slots.default?.());
            },
          }),
        },
      },
    });

  it('should mount successfully', () => {
    const wrapper = mountComponent();
    expect(wrapper.exists()).toBe(true);
  });

  it('should submit create payload with roles', async () => {
    const wrapper = mountComponent();
    formApi.getValues.mockResolvedValueOnce({
      id: '',
      username: 'alice',
      realName: 'Alice',
      password: 'secret',
      roles: ['admin', 'editor'],
      homePath: '',
      avatar: '',
      deptId: null,
      status: 1,
    });

    const buttons = wrapper.findAll('button');
    await buttons[buttons.length - 1].trigger('click');

    expect(createUser).toHaveBeenCalledWith(
      expect.objectContaining({
        username: 'alice',
        roles: '["admin","editor"]',
      }),
    );
  });

  it('should submit password reset for update', async () => {
    const wrapper = mount(Form, {
      props: {
        show: true,
        record: {
          id: 'user-1',
          username: 'alice',
          status: 1,
        },
      },
      global: {
        components: {
          NForm,
        },
        stubs: {
          NButton: defineComponent({
            props: {
              onClick: Function,
            },
            setup(props, { slots }) {
              return () =>
                h(
                  'button',
                  {
                    type: 'button',
                    onClick: props.onClick as any,
                  },
                  slots.default?.(),
                );
            },
          }),
          NModal: defineComponent({
            props: {
              show: Boolean,
            },
            setup(props, { slots }) {
              return () =>
                props.show ? h('div', [slots.default?.(), slots.action?.()]) : null;
            },
          }),
          NFormItem: defineComponent({
            setup(_props, { slots }) {
              return () => h('div', slots.default?.());
            },
          }),
          NInput: defineComponent({
            props: {
              value: String,
              'onUpdate:value': Function,
            },
            setup(props) {
              return () =>
                h('input', {
                  value: props.value,
                  onInput: (event: Event) => {
                    const target = event.target as HTMLInputElement;
                    (props['onUpdate:value'] as any)?.(target.value);
                  },
                });
            },
          }),
          NSpace: defineComponent({
            setup(_props, { slots }) {
              return () => h('div', slots.default?.());
            },
          }),
        },
      },
    });

    await nextTick();
    const buttons = wrapper.findAll('button');
    const changePasswordButton = buttons.find((btn) => {
      const text = btn.text();
      return (
        text.includes('Password') ||
        text.includes('密码') ||
        text.includes('changePassword') ||
        text.includes('common.change')
      );
    });
    if (changePasswordButton) {
      await changePasswordButton.trigger('click');
    } else {
      await buttons[0].trigger('click');
    }

    const input = wrapper.find('input');
    await input.setValue('newpass');

    const actionButtons = wrapper.findAll('button');
    await actionButtons[actionButtons.length - 1].trigger('click');

    expect(updateUser).toHaveBeenCalledWith(
      'user-1',
      expect.objectContaining({
        password: 'newpass',
      }),
    );
  });
});
