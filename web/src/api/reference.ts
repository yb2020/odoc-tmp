import { GetPaperDetailInfoRequest, PaperDetailInfo } from 'go-sea-proto/gen/ts/paper/Paper'
import api from './axios';
import { HEADER_CANCLE_AUTO_ERROR } from './const';
import { SuccessResponse } from './type';

export const getPaperDetailInfo = async (param: GetPaperDetailInfoRequest) => {
  const res = await api.post<SuccessResponse<PaperDetailInfo>>(
    `/paper/getPaperDetailInfo`,
    param,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: true,
      },
    }
  );
  return res.data.data;
};
