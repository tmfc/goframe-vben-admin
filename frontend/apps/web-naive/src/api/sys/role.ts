import { requestClient } from '#/api/request';

export interface Role {
  id: string;
  name: string;
  description: string;
  parentId: string;
  status: number;
}

/**
 * Get role list
 */
export function getRoleList(params?: any) {
  return requestClient.get<any>({ url: '/roles', params });
}

/**
 * Create role
 */
export function createRole(data: Role) {
  return requestClient.post<Role>({ url: '/roles', data });
}

/**
 * Update role
 */
export function updateRole(id: string, data: Role) {
  return requestClient.put<Role>({ url: `/roles/${id}`, data });
}

/**
 * Delete role
 */
export function deleteRole(id: string) {
  return requestClient.delete<void>({ url: `/roles/${id}` });
}

/**
 * Get role
 */
export function getRole(id: string) {
  return requestClient.get<Role>({ url: `/roles/${id}` });
}
