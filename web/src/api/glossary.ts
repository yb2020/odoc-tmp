import api from './axios';
import {
  DeleteGlossaryReq,
  GetGlossaryListReq,
  GetGlossaryListResponse,
  AddGlossaryReq,
  AddGlossaryResponse,
  UpdateGlossaryReq,
} from 'go-sea-proto/gen/ts/translate/GlossaryManage';
import { SuccessResponse } from '@common/api/type';

export const fethGlossrayListByPage = async (params: GetGlossaryListReq) => {
  const res = await api.post<SuccessResponse<GetGlossaryListResponse>>(
    `/glossary/list`,
    params
  );

  return res.data.data;
};

export const deleteGlossaryItems = async (params: DeleteGlossaryReq) => {
  await api.post<SuccessResponse<void>>(
    `/glossary/delete`,
    params
  );
};

export const addGlossaryItem = async (params: AddGlossaryReq) => {
  const res = await api.post<SuccessResponse<AddGlossaryResponse>>(
    `/glossary/add`,
    params
  );

  return res.data.data;
};

export const updateGlossaryItem = async (params: UpdateGlossaryReq) => {
  await api.post<SuccessResponse<void>>(
    `/glossary/update`,
    params
  );
};
