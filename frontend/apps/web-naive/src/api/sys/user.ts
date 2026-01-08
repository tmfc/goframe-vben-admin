import { requestClient } from '#/api/request';
import type { User } from '@vben/types';

/**
 * Get user list
 */
export function getUserList(params?: any) {
  return requestClient.get<any>('/users', { params });
}

/**
 * Create user
 */
export function createUser(data: User) {
  return requestClient.post<User>('/users', data);
}

/**
 * Update user
 */
export function updateUser(id: string, data: User) {
  return requestClient.put<User>(`/users/${id}`, data);
}

/**
 * Delete user
 */
export function deleteUser(id: string) {
  return requestClient.delete<void>(`/users/${id}`);
}

/**
 * Get user
 */
export function getUser(id: string) {
  return requestClient.get<User>(`/users/${id}`);
}
