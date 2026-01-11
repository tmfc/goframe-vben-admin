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

// 批量为指定菜单生成默认按钮（新增/修改/删除）
export function generateButtons(menuId: string) {
  return requestClient.post(`${API_URL}/${menuId}/generate-buttons`);
}
