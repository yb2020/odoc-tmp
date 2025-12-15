import api from './axios';
import { HEADER_CANCLE_AUTO_ERROR, REQUEST_SERVICE_NAME_APP } from './const';
import { SuccessResponse } from './type';
import {
  GroupNoteCommentCreateReq,
  IncrementalGroupCommentReq,
  IncrementalGroupCommentRsp,
  IncrementalGroupNoteReq,
  IncrementalGroupNoteRsp,
  GroupNoteCommentUpdateReq,
  GroupNoteCommentCreateRsp,
  GroupNoteCommentDeleteReq,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/GroupNote';
import qs from 'qs';
import { currentNoteInfo, store } from '@/store';

export const getGroupNotes = async (params: IncrementalGroupNoteReq) => {
  const { data: res } = await api.get<SuccessResponse<IncrementalGroupNoteRsp>>(
    `${REQUEST_SERVICE_NAME_APP}/pdfMark/v2/web/group/fetch/incremental`,
    {
      params,
      paramsSerializer: (params) => {
        const { operationTypes, ...args } = params;

        return qs.stringify(args) + '&' + `operationTypes=${operationTypes.join(',')}`;
      },
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: true,
      },
    }
  );

  return res.data;
};

export const getGroupComments = async (params: IncrementalGroupCommentReq) => {
  const { data: res } = await api.get<SuccessResponse<IncrementalGroupCommentRsp>>(
    `${REQUEST_SERVICE_NAME_APP}/paperNote/group/comment/fetch/incremental`,
    {
      params: {
        ...params,
        noteId: currentNoteInfo.value?.noteId,
        groupId: store.state.base.currentGroupId,
      },
      paramsSerializer: (params) => {
        const { operationTypes, markIds, ...args } = params;

        let str = qs.stringify(args) + '&' + `operationTypes=${operationTypes.join(',')}`;

        if (markIds && markIds.length > 0) {
          str += '&' + `markIds=${markIds.join(',')}`;
        }

        return str;
      },
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: true,
      },
    }
  );

  return res.data;
};

export const addComment = async (params: GroupNoteCommentCreateReq) => {
  const { data: res } = await api.post<SuccessResponse<GroupNoteCommentCreateRsp>>(
    `${REQUEST_SERVICE_NAME_APP}/paperNote/group/comment`,
    {
      ...params,
      noteId: currentNoteInfo.value?.noteId,
    }
  );

  return res;
};

export const updateComment = async (params: GroupNoteCommentUpdateReq) => {
  const { data: res } = await api.put(`${REQUEST_SERVICE_NAME_APP}/paperNote/group/comment`, {
    ...params,
    noteId: currentNoteInfo.value?.noteId,
  });

  return res;
};

export const deleteComment = async (params: GroupNoteCommentDeleteReq) => {
  const { data: res } = await api.delete(`${REQUEST_SERVICE_NAME_APP}/paperNote/group/comment`, {
    data: params,
  });

  return res;
};
