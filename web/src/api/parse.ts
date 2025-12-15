// import {
//   // GetReferenceMarkersRequest,
//   // GetReferenceMarkersResponse,
//   // GetFiguresAndTablesListRequest,
//   // GetFiguresAndTablesListResponse,
//   // PdfFigureAndTableInfo,
//   // ReferenceMarker,
//   // FigureAndTableReferenceMarker,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse';
import { 
  GetReferenceMarkersRequest, 
  GetReferenceMarkersResponse, 
  ReferenceMarker,
  GetFiguresAndTablesListRequest, 
  GetFiguresAndTablesListResponse, 
  FigureAndTableReferenceMarker,
  PdfFigureAndTableInfo,
  GetCatalogueRequest
 } from 'go-sea-proto/gen/ts/pdf/PdfParse';
import i18n from '../locals/i18n';
import api from './axios'
import { HEADER_CANCLE_AUTO_ERROR } from './const';
import { ResponseError, SuccessResponse } from './type';

export const getReferenceMarkers = async (params: GetReferenceMarkersRequest) => {
  const { data: res } = await api.post<SuccessResponse<GetReferenceMarkersResponse>>(
    '/pdf/parser/getReferenceMarkers',
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: true,
      },
    }
  );

  return {
    needFetch: res.data.needFetch,
    referenceMarkers: res.data.markers || [],
    figureAndTableMarkers: res.data.figureAndTableMarkers || [],
  };
};

export const firstReferenceMarkersFetch = async (params: GetReferenceMarkersRequest) => {
  const res = await getReferenceMarkers(params);
  if (!res.needFetch) {
    return {
      referenceMarkers: res.referenceMarkers || [],
      figureAndTableMarkers: res.figureAndTableMarkers || [],
    };
  }

  return null;
};

export const getReferenceMarkersFinal = (
  params: GetReferenceMarkersRequest
): Promise<{ referenceMarkers: ReferenceMarker[]; figureAndTableMarkers: FigureAndTableReferenceMarker[] }> => {
  // 每次调用都创建独立的 count 和 error 变量
  let count = 100;
  let error: null | ResponseError = null;
  
  return new Promise(async (resolve, reject) => {
    const fetch = async () => {
      try {
        const data = await getReferenceMarkers(params);
        
        if (!data.needFetch) {
          return resolve({
            referenceMarkers: data.referenceMarkers || [],
            figureAndTableMarkers: data.figureAndTableMarkers || [],
          });
        }
      } catch (err) {
        error = err as ResponseError;
      }
      
      if (count-- <= 0) {
        return reject(error || Error(i18n.global.t('message.resourceSchedulingTimeoutTip')));
      }
      
      setTimeout(() => {
        fetch();
      }, 1000);
    };
    fetch();
  });
};

export const getFiguresAndTables = async (params: GetFiguresAndTablesListRequest) => {
  const res = await api.post<SuccessResponse<GetFiguresAndTablesListResponse>>(
    `/pdf/parser/getFiguresAndTables`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: true,
      },
    }
  );
  return res.data.data;
};

export const firstFigureAndTableFetch = async (params: GetFiguresAndTablesListRequest) => {
  const res = await getFiguresAndTables(params);

  if (!res.needFetch) {
    const list = res.figureAndTableList || [];
    return {
      list,
      total: res.pageResp?.total || list.length,
    };
  }
  return null;
};

export const getFiguresAndTablesFinal = (
  params: GetFiguresAndTablesListRequest
): Promise<{ list: PdfFigureAndTableInfo[]; total: number }> => {
  // 每次调用都创建独立的 count 和 error 变量
  let count = 100;
  let error: null | ResponseError = null;
  
  return new Promise(async (resolve, reject) => {
    const fetch = async () => {
      try {
        const data = await getFiguresAndTables(params);
        if (!data.needFetch) {
          const list = data.figureAndTableList || [];
          return resolve({
            list,
            total: data.pageResp?.total || list.length,
          });
        }
      } catch (err) {
        error = err as ResponseError;
      }
      if (count-- <= 0) {
        return reject(error || Error(i18n.global.t('message.resourceSchedulingTimeoutTip')));
      }
      setTimeout(() => {
        fetch();
      }, 1000);
    };

    fetch();
  });
};

/**
 * 重新解析文档
 * @param params - 请求参数
 */
export const reParsePaper = async (params: GetCatalogueRequest): Promise<void> => {
  await api.post('/pdf/parser/reParse', params)
}