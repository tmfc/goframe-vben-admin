import { describe, it, expect, vi } from 'vitest';
import {
  getRoleList,
  createRole,
  updateRole,
  deleteRole,
  getRole,
  getRolePermissions,
  assignRolePermissions,
} from '../role';
import { requestClient } from '#/api/request';

vi.mock('#/api/request', () => ({
  requestClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('Role API', () => {
  it('should call getRoleList with correct parameters', async () => {
    await getRoleList({ page: 1, pageSize: 10 });
    expect(requestClient.get).toHaveBeenCalledWith('/sys-role', {
      params: { page: 1, pageSize: 10 },
    });
  });

  it('should call createRole with correct parameters', async () => {
    const data = { name: 'test' };
    // @ts-ignore
    await createRole(data);
    expect(requestClient.post).toHaveBeenCalledWith('/sys-role', data);
  });

  it('should call updateRole with correct parameters', async () => {
    const id = '123';
    const data = { name: 'test' };
    // @ts-ignore
    await updateRole(id, data);
    expect(requestClient.put).toHaveBeenCalledWith(`/sys-role/123`, data);
  });

  it('should call deleteRole with correct parameters', async () => {
    const id = '123';
    await deleteRole(id);
    expect(requestClient.delete).toHaveBeenCalledWith(`/sys-role/123`);
  });

  it('should call getRole with correct parameters', async () => {
    const id = '123';
    await getRole(id);
    expect(requestClient.get).toHaveBeenCalledWith(`/sys-role/123`);
  });

  it('should call getRolePermissions with correct parameters', async () => {
    const id = '123';
    await getRolePermissions(id);
    expect(requestClient.get).toHaveBeenCalledWith(`/sys-role/123/permissions`);
  });

  it('should call assignRolePermissions with correct parameters', async () => {
    const id = '123';
    const data = { permissionIds: [1, 2] };
    await assignRolePermissions(id, data);
    expect(requestClient.post).toHaveBeenCalledWith(
      `/sys-role/123/assign-permissions`,
      data,
    );
  });
});
