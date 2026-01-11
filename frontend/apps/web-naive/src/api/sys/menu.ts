import { requestClient } from '#/api/request';

const API_URL = '/sys-menu';

export interface MenuItem {
  id: string;
  parentId?: string;
  name: string;
  path: string;
  component?: string;
  icon?: string;
  type: string; // '0' | '1' | '2' (Directory, Menu, Button)
  status: number;
  order: number;
  permissionCode?: string;
  visible?: number;
  meta?: string; // JSON string in DB, might need parsing or comes as string
  tenantId?: string;
  creatorId?: number;
  modifierId?: number;
  deptId?: number;
  createdAt?: string;
  updatedAt?: string;
  children?: MenuItem[];
}

export interface CreateMenuParams {
  name: string;
  path?: string;
  component?: string;
  icon?: string;
  type: string;
  parentId?: string;
  meta?: string;
  status?: number;
  order?: number;
  permissionCode?: string;
}

export interface UpdateMenuParams extends CreateMenuParams {
  id: string;
}

export interface GetMenuListParams {
  name?: string;
  status?: string;
}

export interface GenerateButtonsResult {
  created: number;
  skipped: number;
}

export function getMenuList(params?: GetMenuListParams) {
  return requestClient.get<{ list: MenuItem[]; total: number }>(`${API_URL}/list`, { params });
}

export function createMenu(data: CreateMenuParams) {
  return requestClient.post<{ id: string }>(API_URL, data);
}

export function updateMenu(id: string, data: UpdateMenuParams) {
  return requestClient.put<void>(`${API_URL}/${id}`, data);
}

export function deleteMenu(id: string) {
  return requestClient.delete<void>(`${API_URL}/${id}`);
}

// 批量为指定菜单生成默认按钮（新增/修改/删除）
export function generateButtons(menuId: string) {
  return requestClient.post<GenerateButtonsResult>(`${API_URL}/${menuId}/generate-buttons`);
}
