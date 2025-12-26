import { requestClient } from '#/api/request';

const API_URL = '/sys-menu';

export function getMenuList(params) {
  return requestClient.get({ url: `${API_URL}/list`, params });
}

export function createMenu(data) {
  return requestClient.post({ url: API_URL, data });
}

export function updateMenu(id, data) {
  return requestClient.put({ url: `${API_URL}/${id}`, data });
}

export function deleteMenu(id) {
  return requestClient.delete({ url: `${API_URL}/${id}` });
}
