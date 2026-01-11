import { describe, it, expect, vi } from 'vitest';
import {
  getMenuList,
  createMenu,
  updateMenu,
  deleteMenu,
  generateButtons,
} from '../menu';
import { requestClient } from '#/api/request';

vi.mock('#/api/request', () => ({
  requestClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('Menu API', () => {
  it('should call getMenuList with correct parameters', async () => {
    const params = { name: 'menu1', status: '1' };
    await getMenuList(params);
    expect(requestClient.get).toHaveBeenCalledWith('/sys-menu/list', {
      params,
    });
  });

  it('should call createMenu with correct parameters', async () => {
    const data = { name: 'menu1', type: 'menu' };
    await createMenu(data);
    expect(requestClient.post).toHaveBeenCalledWith('/sys-menu', data);
  });

  it('should call updateMenu with correct parameters', async () => {
    const data = { name: 'menu1', type: 'menu', id: '123' };
    await updateMenu('123', data);
    expect(requestClient.put).toHaveBeenCalledWith('/sys-menu/123', data);
  });

  it('should call deleteMenu with correct parameters', async () => {
    await deleteMenu('123');
    expect(requestClient.delete).toHaveBeenCalledWith('/sys-menu/123');
  });

  it('should call generateButtons with correct parameters', async () => {
    await generateButtons('123');
    expect(requestClient.post).toHaveBeenCalledWith('/sys-menu/123/generate-buttons');
  });
});
