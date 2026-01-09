import { describe, it, expect, vi } from 'vitest';
import {
  getPermissionList,
  createPermission,
  updatePermission,
  deletePermission,
  getPermission,
} from '../permission';
import { requestClient } from '#/api/request';

vi.mock('#/api/request', () => ({
  requestClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('Permission API', () => {
  it('should call getPermissionList with correct parameters', async () => {
    await getPermissionList({ page: 1, pageSize: 100 });
    expect(requestClient.get).toHaveBeenCalledWith('/sys-permission/list', {
      params: { page: 1, pageSize: 100 },
    });
  });

  it('should call createPermission with correct parameters', async () => {
    const data = { name: 'perm-a', status: 1 };
    await createPermission(data);
    expect(requestClient.post).toHaveBeenCalledWith('/sys-permission', data);
  });

  it('should call updatePermission with correct parameters', async () => {
    const data = { name: 'perm-a' };
    await updatePermission(123, data);
    expect(requestClient.put).toHaveBeenCalledWith('/sys-permission/123', data);
  });

  it('should call deletePermission with correct parameters', async () => {
    await deletePermission(123);
    expect(requestClient.delete).toHaveBeenCalledWith('/sys-permission/123');
  });

  it('should call getPermission with correct parameters', async () => {
    await getPermission(123);
    expect(requestClient.get).toHaveBeenCalledWith('/sys-permission/123');
  });
});
