import { SuccessResponse } from './type';
import { $PaperNoteItem } from './note';
import { $UserInfoData } from './user';
import { REQUEST_SERVICE_NAME_APP } from './const';
import { api } from './axios';

export type $AuthorInfo = $UserInfoData;

type $GetHotPaperNoteData = {
  list: $PaperNoteItem[];
  total: number;
};

interface $GetHotPaperNoteListReq {
  currentPage: number;
  pageSize: number;
}

type $GetHotPaperNoteRsp = SuccessResponse<$GetHotPaperNoteData>;

export const fetchHotPaperNoteList = async (
  params: $GetHotPaperNoteListReq
) => {
  const res = await api.post<$GetHotPaperNoteRsp>(
    `${REQUEST_SERVICE_NAME_APP}/aiKnowledge/hotNote/list`,
    {
      ...params,
      orderType: 1,
      order: 0,
    }
  );
  return res.data.data;
};
