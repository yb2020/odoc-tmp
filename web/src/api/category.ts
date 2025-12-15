/*
 * Created Date: May 14th 2021, 3:56:09 pm
 * Author: zhoupengcheng
 * -----
 * Last Modified: August 12th 2021, 3:55:00 pm
 */
import api from './axios'
import { HEADER_CANCLE_AUTO_ERROR } from './const'
import { ResponseError, SuccessResponse } from './type'
// import {
//   GetCatalogueRequest,
//   GetCatalogueResponse,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse'
import {GetCatalogueRequest, GetCatalogueResponse} from 'go-sea-proto/gen/ts/pdf/PdfParse'

export const getCatalog = (() => {
  let count = 4
  let error: null | ResponseError = null
  return (params: GetCatalogueRequest): Promise<GetCatalogueResponse> => {
    return new Promise(async (resolve, reject) => {
      const fetch = async () => {
        try {
          const { data: res } = await api.post<
            SuccessResponse<GetCatalogueResponse>
          >(`/pdf/parser/getCatalogue`, params, {
            headers: {
              [HEADER_CANCLE_AUTO_ERROR]: true,
            },
          })
          if (!res.data.needFetch) {
            resolve(res.data)
            return
          }
        } catch (err) {
          error = err as ResponseError
        }

        if (count-- <= 0) {
          reject(error)
          return
        }
        setTimeout(() => {
          fetch()
        }, 2000)
      }

      fetch()
    })
  }
})()
