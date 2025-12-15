import api from './axios';
import { REQUEST_SERVICE_NAME_APP } from './const';
import { SuccessResponse } from './type';
import {
  GroupNoteReq,
  GroupNoteRsp,
  GroupDocReq,
  GroupDocRsp
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/GroupNote';

export type $GroupProceed = {
  id: string;
  name: string;
  count: number;
  modifyDate: number;
  creator: string;
  creatorId: string;
};
// 获取创建小组列表
interface $GetMyCreatedGroupListReq {
  currentPage?: number;
  pageSize?: number;
}
type $GetMyCreatedGroupListRsp = SuccessResponse<{
  total: number;
  list: $GroupProceed[];
}>;

export const getMyCreatedGroupList = async (params: $GetMyCreatedGroupListReq) => {
  try {
    const res = await api.post<$GetMyCreatedGroupListRsp>(
      `${REQUEST_SERVICE_NAME_APP}/group/getMyCreatedGroupList`,
      params
    );
    return res.data.data;
  } catch (error) {
    return {
      total: 0,
      list: [],
    }
  }
};

interface $GetMyJoinGroupListReq {
  currentPage?: number;
  pageSize?: number;
}
type $GetMyJoinGroupListRsp = SuccessResponse<{
  total: number;
  list: $GroupProceed[];
}>;

export const getMyJoinGroupList = async (params: $GetMyJoinGroupListReq) => {
  try {
    const res = await api.post<$GetMyJoinGroupListRsp>(
      `${REQUEST_SERVICE_NAME_APP}/group/getMyJoinGroupList`,
      params
    );
    return res.data.data;
  } catch (error) {
    return {
      total: 0,
      list: [],
    }
  }
  
};

export const getGroupNote = async (params: GroupNoteReq) => {
  const res = await api.post<SuccessResponse<GroupNoteRsp>>(
    `${REQUEST_SERVICE_NAME_APP}/paperNote/group/note`,
    params
  );  
  return res.data.data;
};

export const addPaperToGroup = async (params: GroupDocReq) => {
  const res = await api.post<SuccessResponse<GroupDocRsp>>(
    `${REQUEST_SERVICE_NAME_APP}/paperNote/group/addPaperToGroup`,
    params
  );  
  return res.data.data;
};
