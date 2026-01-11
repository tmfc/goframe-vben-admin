import { requestClient } from '#/api/request';

export interface TenantModel {
  id?: string | number;
  name: string;
  code: string;
  status: number;
  contactName?: string;
  contactMobile?: string;
  startAt?: string;
  expireAt?: string;
  packageVersion?: string;
  domain?: string;
  createdAt?: string;
  updatedAt?: string;
}

/**
 * Get tenant list
 */
export function getTenantList(params?: any) {
  return requestClient.get<any>('/sys-tenant', { params });
}

/**
 * Create tenant
 */
export function createTenant(data: TenantModel) {
  return requestClient.post<TenantModel>('/sys-tenant', data);
}

/**
 * Update tenant
 */
export function updateTenant(id: string | number, data: TenantModel) {
  return requestClient.put<TenantModel>(`/sys-tenant/${id}`, data);
}

/**
 * Delete tenant
 */
export function deleteTenant(id: string | number) {
  return requestClient.delete<void>(`/sys-tenant/${id}`);
}

/**
 * Get tenant
 */
export function getTenant(id: string | number) {
  return requestClient.get<TenantModel>(`/sys-tenant/${id}`);
}
