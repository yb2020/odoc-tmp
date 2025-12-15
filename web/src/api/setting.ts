import {
  MarkFinishReq,
  // UpdateReadStatusRequest,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/doc/UserDocManage';
import { UpdateReadStatusRequest } from 'go-sea-proto/gen/ts/doc/UserDocManage';
import {
  // HEADER_CANCLE_AUTO_ERROR,
  // REQUEST_SERVICE_NAME_APP,
  REQUEST_SERVICE_NAME_DOC,
} from './const';

import api from './axios';
import { LeftSideBarType } from '../components/Left/type';
import { RightSideBarType } from '../components/Right/TabPanel/type';
import { SuccessResponse } from './type';
import { isInElectron } from '@/util/env';
import { LangType } from '../stores/copilotType';

// 0：代表浏览器
// 1：代表win客户端
// 2：代表ios客户端
const getClientType = () => {
  return isInElectron() ? 1 : 0;
};

export type UserSettingInfo = Partial<{
  scale: number;
  scalePresetValue: string;
  rightTab: RightSideBarType;
  rightSubTab?: number | string;
  rightWidth: number;
  rightShow: boolean;
  sideBarTab: LeftSideBarType;
  sideBarWidth: number;
  sideBarShow: boolean;
  rightTabBars: { key: RightSideBarType; shown: boolean }[];
  copilotLanguage: LangType;
  toolBarVisible: boolean;
  toolBarHeadVisible: boolean;
  toolBarNoteVisible: boolean;
}>;

export const getSetting = async () => {
  try {
    const { data: res } = await api.post<SuccessResponse<{ setting: string }>>(
      `/pdf/pdfReader/getSetting`,
      { clientType: getClientType() },
    );
    return JSON.parse(res.data.setting) as UserSettingInfo;
  } catch (error) {
    return {};
  }
};

export const setSetting = async (param: UserSettingInfo) => {
  try {
    await api.post(
      `/pdf/pdfReader/recordSetting`,
      {
        clientType: getClientType(),
        setting: JSON.stringify(param),
      },
    );
  } catch (error) {}
};

interface GetLocationRequest {
  noteId: string;
}

export interface NoteLocationInfo extends UserSettingInfo {
  currentPage?: number;
  currentGroupId?: string;
}

export const getLocation = async (param: GetLocationRequest) => {
  try {
    const { data: res } = await api.post<SuccessResponse<{ location: string }>>(
      `/note/noteReadLocation/getLocation`,
      param,
    );

    return JSON.parse(res.data.location) as NoteLocationInfo;
  } catch (error) {
    return {};
  }
};

export const setLocation = async (
  params: { location: NoteLocationInfo } & GetLocationRequest
) => {
  try {
    await api.post<SuccessResponse<null>>(
      `/note/noteReadLocation/record`,
      { location: JSON.stringify(params.location), noteId: params.noteId },
    );
  } catch (error) {}
};

/**
 * @deprecated
 */
export const markReadingAsRead = async (params: MarkFinishReq) => {
  try {
    await api.put<SuccessResponse<null>>(
      `${REQUEST_SERVICE_NAME_DOC}/userDoc/mark/finish`,
      params,
    );
  } catch (error) {}
};

export const updateDocReadStatus = async (params: UpdateReadStatusRequest) => {
  const res = await api.post<object>(
    `/doc/userDoc/updateReadStatus`,
    params
  );

  return res.data;
};
