import { requestClient } from '#/api/request';

const API_URL = '/sys-menu';

export function getMenuList(params) {
  return requestClient.get(`${API_URL}/list`, { params });
}

export function createMenu(data) {
  return requestClient.post(API_URL, data);
}

export function updateMenu(id, data) {
  return requestClient.put(`${API_URL}/${id}`, data);
}

export function deleteMenu(id) {
  return requestClient.delete(`${API_URL}/${id}`);
}
