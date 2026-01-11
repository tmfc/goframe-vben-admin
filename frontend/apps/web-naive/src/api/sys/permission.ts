import { requestClient } from '#/api/request';

export interface PermissionItem {
  id: number | string;
  name: string;
  parentId?: number | string;
  status?: number;
  description?: string;
  createdAt?: string;
  updatedAt?: string;
  children?: PermissionItem[];
  // Include other fields from entity if needed
  deptId?: number | string;
  creatorId?: number | string;
  modifierId?: number | string;
}

export interface CreatePermissionParams {
  name: string;
  description?: string;
  parent_id?: number | string;
  status?: number;
  dept_id?: number | string;
}

export interface UpdatePermissionParams extends CreatePermissionParams {
  id: number | string;
}

export interface GetPermissionListParams {
  page?: number;
  pageSize?: number;
  name?: string;
  status?: string;
}

export interface PermissionListResult {
  items: PermissionItem[];
  total: number;
}

/**
 * Get permission list
 */
export function getPermissionList(params?: GetPermissionListParams) {
  return requestClient.get<PermissionListResult>('/sys-permission/list', { params });
}

/**
 * Create permission
 */
export function createPermission(data: CreatePermissionParams) {
  return requestClient.post<{ id: number }>('/sys-permission', data);
}

/**
 * Update permission
 */
export function updatePermission(id: string | number, data: CreatePermissionParams) {
  return requestClient.put<void>(`/sys-permission/${id}`, data);
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
  return requestClient.get<PermissionItem>(`/sys-permission/${id}`);
}

/**
 * Get permission tree
 */
export function getPermissionTree() {
  return requestClient.get<{ list: PermissionItem[] }>('/sys-permission/tree');
}
