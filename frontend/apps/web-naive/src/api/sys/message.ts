import { request } from '#/api/request';

export namespace MessageApi {
  export interface MessageItem {
    id: number;
    title: string;
    content: string;
    type: string;
    created_at: string;
    read_at: string | null;
    jump_path: string | null;
    jump_params: string | null;
  }

  export interface GetListParams {
    page: number;
    size: number;
  }

  export interface GetListResult {
    list: MessageItem[];
    total: number;
  }
}

/**
 * 获取我的消息列表
 */
export function getMessageListApi(params: MessageApi.GetListParams) {
  return request.get<MessageApi.GetListResult>('/message/list', { params });
}

/**
 * 标记消息为已读
 */
export function setMessageReadApi(ids: number[]) {
  return request.post('/message/read', { ids });
}

/**
 * 全部标记为已读
 */
export function setAllMessageReadApi() {
  return request.post('/message/read-all');
}
