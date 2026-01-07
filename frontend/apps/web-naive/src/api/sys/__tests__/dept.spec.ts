import { describe, it, expect, vi } from 'vitest';
import { getDeptList, createDept, updateDept, deleteDept, getDept } from '../dept';
import { requestClient } from '#/api/request';

vi.mock('#/api/request', () => ({
  requestClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('Department API', () => {
  it('should call getDeptList with correct parameters', async () => {
    await getDeptList({ page: 1, pageSize: 10 });
    expect(requestClient.get).toHaveBeenCalledWith({ url: '/sys-dept/list', params: { page: 1, pageSize: 10 } });
  });

  it('should call createDept with correct parameters', async () => {
    const data = { name: 'Test Dept' };
    await createDept(data);
    expect(requestClient.post).toHaveBeenCalledWith({ url: '/sys-dept', data });
  });

  it('should call updateDept with correct parameters', async () => {
    const id = '123';
    const data = { name: 'Updated Dept' };
    await updateDept(id, data);
    expect(requestClient.put).toHaveBeenCalledWith({ url: '/sys-dept/123', data });
  });

  it('should call deleteDept with correct parameters', async () => {
    const id = '123';
    await deleteDept(id);
    expect(requestClient.delete).toHaveBeenCalledWith({ url: '/sys-dept/123' });
  });

  it('should call getDept with correct parameters', async () => {
    const id = '123';
    await getDept(id);
    expect(requestClient.get).toHaveBeenCalledWith({ url: '/sys-dept/123' });
  });
});
