import { requestClient } from '#/api/request';

const API_URL = '/sys-dept';

/**
 * Get department list
 */
export function getDeptList(params?: any) {
  return requestClient.get({ url: `${API_URL}/list`, params });
}

/**
 * Get department tree
 */
export function getDeptTree(params?: any) {
  return requestClient.get({ url: `${API_URL}/tree`, params });
}

/**
 * Create a new department
 */
export function createDept(data: any) {
  return requestClient.post({ url: API_URL, data });
}

/**
 * Get department details
 */
export function getDept(id: string) {
  return requestClient.get({ url: `${API_URL}/${id}` });
}

/**
 * Update department
 */
export function updateDept(id: string, data: any) {
  return requestClient.put({ url: `${API_URL}/${id}`, data });
}

/**
 * Delete department
 */
export function deleteDept(id: string) {
  return requestClient.delete({ url: `${API_URL}/${id}` });
}
