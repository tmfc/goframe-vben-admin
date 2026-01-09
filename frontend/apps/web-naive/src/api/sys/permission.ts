import { requestClient } from '#/api/request';

export interface Permission {
  id: number | string;
  name: string;
  parentId?: number | string;
  status?: number;
  description?: string;
}

export interface CreatePermissionPayload {
  name: string;
  description?: string;
  parent_id?: number;
  status?: number;
}

/**
 * Get permission list
 */
export function getPermissionList(params?: any) {
  return requestClient.get<any>('/sys-permission/list', { params });
}

/**
 * Create permission
 */
export function createPermission(data: CreatePermissionPayload) {
  return requestClient.post('/sys-permission', data);
}

/**
 * Update permission
 */
export function updatePermission(id: string | number, data: CreatePermissionPayload) {
  return requestClient.put(`/sys-permission/${id}`, data);
}

/**
 * Delete permission
 */
export function deletePermission(id: string | number) {
  return requestClient.delete<void>(`/sys-permission/${id}`);
}

/**
 * Get permission
 */
export function getPermission(id: string | number) {
  return requestClient.get<Permission>(`/sys-permission/${id}`);
}
