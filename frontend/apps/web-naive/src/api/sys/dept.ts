import { requestClient } from '#/api/request';

const API_URL = '/sys-dept';

export interface Dept {
  id: string;
  name: string;
  order: number;
  status: number;
}

/**
 * Get department list
 */
export function getDeptList(params?: any) {
  return requestClient.get<{ list: Dept[]; total: number }>(
    `${API_URL}/list`,
    { params },
  );
}

/**
 * Get department tree
 */
export function getDeptTree(params?: any) {
  return requestClient.get<{ list: Dept[] }>(`${API_URL}/tree`, { params });
}

/**
 * Create a new department
 */
export function createDept(data: any) {
  return requestClient.post(API_URL, data);
}

/**
 * Get department details
 */
export function getDept(id: string) {
  return requestClient.get<Dept>(`${API_URL}/${id}`);
}

/**
 * Update department
 */
export function updateDept(id: string, data: any) {
  return requestClient.put(`${API_URL}/${id}`, data);
}

/**
 * Delete department
 */
export function deleteDept(id: string) {
  return requestClient.delete(`${API_URL}/${id}`);
}
