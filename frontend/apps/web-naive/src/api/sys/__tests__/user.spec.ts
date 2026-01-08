import { describe, it, expect, vi } from 'vitest';
import { getUserList, createUser, updateUser, deleteUser, getUser } from '../user';
import { requestClient } from '#/api/request';

vi.mock('#/api/request', () => ({
  requestClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('User API', () => {
  it('should call getUserList with correct parameters', async () => {
    await getUserList({ page: 1, pageSize: 10 });
    expect(requestClient.get).toHaveBeenCalledWith('/users', {
      params: { page: 1, pageSize: 10 },
    });
  });

  it('should call createUser with correct parameters', async () => {
    const data = { username: 'test' };
    // @ts-ignore
    await createUser(data);
    expect(requestClient.post).toHaveBeenCalledWith('/users', data);
  });

  it('should call updateUser with correct parameters', async () => {
    const id = '123';
    const data = { username: 'test' };
    // @ts-ignore
    await updateUser(id, data);
    expect(requestClient.put).toHaveBeenCalledWith('/users/123', data);
  });

  it('should call deleteUser with correct parameters', async () => {
    const id = '123';
    await deleteUser(id);
    expect(requestClient.delete).toHaveBeenCalledWith('/users/123');
  });

  it('should call getUser with correct parameters', async () => {
    const id = '123';
    await getUser(id);
    expect(requestClient.get).toHaveBeenCalledWith('/users/123');
  });
});
