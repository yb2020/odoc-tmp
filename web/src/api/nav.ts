
import api from './axios';
import { SuccessResponse } from './type';


import {
    CreateWebsiteRequest,
    CreateWebsiteResponse,
    DeleteWebsiteRequest,
    GetWebsiteByIdRequest,
    GetWebsiteByIdResponse,
    UpdateWebsiteRequest,
    UpdateWebsiteResponse,
    GetWebsiteListRequest,
    GetWebsiteListResponse,
    ReorderWebsitesRequest,
    ReorderWebsitesResponse,
    } from 'go-sea-proto/gen/ts/nav/Website';




// @api_path: /api/nav/website/create
// @method: POST
// @summary: 创建用户学术网站
export const createWebsite = async (params: CreateWebsiteRequest) => {
    const { data: res } = await api.post<SuccessResponse<CreateWebsiteResponse>>(
        `/nav/website/create`,
        params
      );
    return res.data;
};

// @api_path: /api/nav/website/delete
// @method: POST
// @summary: 删除用户学术网站
export const deleteWebsite = async (params: DeleteWebsiteRequest) => {
    const { data: res } = await api.post<SuccessResponse<null>>(
        `/nav/website/delete`,
        params
      );
    return res.data;
};  

// @api_path: /api/nav/website/getById
// @method: GET
// @summary: 获取用户学术网站
export const getWebsiteById = async (params: GetWebsiteByIdRequest) => {
    const { data: res } = await api.get<SuccessResponse<GetWebsiteByIdResponse>>(
        `/nav/website/getById`,
        {
            params
        }
      );
    return res.data;
};

// @api_path: /api/nav/website/update
// @method: POST
// @summary: 更新用户学术网站
export const updateWebsite = async (params: UpdateWebsiteRequest) => {
    const { data: res } = await api.post<SuccessResponse<UpdateWebsiteResponse>>(
        `/nav/website/update`,
        params
      );
    return res.data;
};
        

// @api_path: /api/nav/website/getList
// @method: GET
// @summary: 获取用户学术网站列表
export const getWebsiteList = async (params: GetWebsiteListRequest) => {
    const { data: res } = await api.get<SuccessResponse<GetWebsiteListResponse>>(
        `/nav/website/getList`,
        {
            params
        }
      );
    return res.data;
};
        
// @api_path: /api/nav/website/reorder
// @method: POST
// @summary: 重新排序用户学术网站
export const reorderWebsite = async (params: ReorderWebsitesRequest) => {
    const { data: res } = await api.post<SuccessResponse<ReorderWebsitesResponse>>(
        `/nav/website/reorder`,
        params
      );
    return res.data;
};
        