import { requestClient } from '#/api/request';

export interface Role {
  id: string;
  name: string;
  description: string;
  parentId: string;
  deptId?: string;
  status: number;
}

/**
 * Get role list
 */
export function getRoleList(params?: any) {
  return requestClient.get<any>('/sys-role', { params });
}

/**
 * Create role
 */
export function createRole(data: Role) {
  return requestClient.post<Role>('/sys-role', data);
}

/**
 * Update role
 */
export function updateRole(id: string, data: Role) {
  return requestClient.put<Role>(`/sys-role/${id}`, data);
}

/**
 * Delete role
 */
export function deleteRole(id: string) {
  return requestClient.delete<void>(`/sys-role/${id}`);
}

/**
 * Get role
 */
export function getRole(id: string) {
  return requestClient.get<Role>(`/sys-role/${id}`);
}
