import {
  SortAnnotationRequest,
  SortAnnotationResponse,
  // WebNoteAnnotationModel,
  // HotSelectRequest,
  // WebDrawV2,
  BatchDeleteDrawReq,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Web';
import { WebNoteAnnotationModel, WebDrawV2 } from 'go-sea-proto/gen/ts/note/web'

// import {
//   AddTagToAnnotateReq,
//   CreateAnnotateTagReq,
//   CreateAnnotateTagResp,
//   DeleteTagToAnnotateReq,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Tag';

import {AddTagToAnnotateRequest, CreateAnnotateTagRequest, CreateAnnotateTagResponse, DeleteAnnotateTagRequest, DeleteTagToAnnotateRequest, RenameAnnotateTagRequest, RenameAnnotateTagResponse} from 'go-sea-proto/gen/ts/note/PaperNoteAnnotateTag';
// import {
//   // DeleteShapeReq,
//   // SaveShapeReq,
//   // SaveShapeResponse,
//   // UpdateShapeReq,
// } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Shape';
import { IDEAAnnotateType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { GetNoteShapesByNoteIdRequest, GetNoteShapesByNoteIdResponse,
  UpdateShapeRequest,
  SaveShapeRequest,
  SaveShapeResponse,
  DeleteShapeRequest, } from 'go-sea-proto/gen/ts/note/NoteShape'
import { GetNoteAnnotationListByNoteIdResponse, GetDrawNoteAnnotationListByNoteIdResponse, HotSelectRequest, HotSelectResponse } from 'go-sea-proto/gen/ts/pdf/PdfMark'
// import type { IDEAAnnotateType } from 'go-sea-proto/gen/ts/common/IDEAAnnotateType'
import {
  getIpadUnit,
  // getTextBody,
  PageShapeTextMap,
  TextAnnotation,
} from '@idea/pdf-annotate-core';
import api from './axios';
import {
  HEADER_CANCLE_AUTO_ERROR,
  REQUEST_SERVICE_NAME_APP,
  REQUEST_SERVICE_NAME_NOTE,
} from './const';
import { SuccessResponse } from './type';
import { usePdfStore } from '../stores/pdfStore';
import { currentNoteInfo } from '../store';

const CONVERTED = 'converted' + Math.random();

export const getByNote = async (params: {
  noteId: string;
  groupId?: string;
}) => {
  const { data: res } = await api.get<SuccessResponse<GetNoteAnnotationListByNoteIdResponse>>(
    `/pdf/pdfMark/v3/web/getByNote`, 
  {
    params: {
      ...params,
      handtext: true,
    },
  });

  const annotationList: WebNoteAnnotationModel[] = [];
  const handTextList: TextAnnotation[] = [];

  res.data.annotations.forEach((item) => {
    if (item.type === IDEAAnnotateType.IDEAAnnotateTypeTextBox) {
      // if (item.textBox?.id && !getTextBody(item)?.content) {
      //   deleteAnnotation({
      //     annotationId: item.textBox.id,
      //     noteId: params.noteId,
      //     groupId: params.groupId,
      //   });

      //   return;
      // }

      handTextList.push(item);
      return;
    }

    if (
      ![
        IDEAAnnotateType.IDEAAnnotateTypeComment,
        IDEAAnnotateType.IDEAAnnotateTypeRect,
        IDEAAnnotateType.IDEAAnnotateTypeHotRect,
      ].includes(item.type)
    ) {
      return;
    }

    if (!item.tags) {
      item.tags = [];
    }

    if (item.select) {
      if (item.select.idea === undefined) {
        item.select.idea = '';
      }
    }

    if (item.rect) {
      if (item.rect.idea === undefined) {
        item.rect.idea = '';
      }
    }

    annotationList.push(item);
  });

  return {
    annotationList,
    handTextList,
  };
};

interface HandwritePageParams {
  fetchAll: boolean;
  pageNumbers?: number[];
}

export const getHandwrite = async (
  params: {
    noteId: string;
    groupId: string;
  },
  pageNumbers?: number[]
) => {
  const pageParams: HandwritePageParams = pageNumbers?.length
    ? {
        fetchAll: false,
        pageNumbers,
      }
    : { fetchAll: true };

  const { data: res } = await api.get<
    SuccessResponse<GetDrawNoteAnnotationListByNoteIdResponse>
  >(`/pdf/pdfMark/v3/web/draw/getByNote`, {
    params: {
      ...params,
      ...pageParams,
      handwrite: true,
    },
    timeout: 20000,
  });

  return res.data.annotations.map((item) => item.webDrawV2).filter(Boolean) as WebDrawV2[];
};

export const getShapeList = async (params: GetNoteShapesByNoteIdRequest) => {
  const { data: res } = await api.post<SuccessResponse<GetNoteShapesByNoteIdResponse>>(
    `/note/noteShape/getList`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );

  return res.data.list || [];
};

export const saveShape = async (params: SaveShapeRequest) => {
  const { data: res } = await api.post<SuccessResponse<SaveShapeResponse>>(
    `/note/noteShape/save`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );

  return res.data.shapeId;
};

export const updateShape = async (params: UpdateShapeRequest) => {
  const { data: res } = await api.post<SuccessResponse<void>>(
    `/note/noteShape/update`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );

  return res.data;
};

export const deleteShape = async (params: DeleteShapeRequest) => {
  const { data: res } = await api.post<SuccessResponse<void>>(
    `/note/noteShape/delete`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );

  return res.data;
};

export const deleteHandwrite = async (params: BatchDeleteDrawReq) => {
  const { data } = await api.delete(
    `${REQUEST_SERVICE_NAME_APP}/pdfMark/v2/web/pencil/batch/delete`,
    {
      data: params,
    }
  );

  return data;
};

export const getHotAnnotations = async (params: HotSelectRequest) => {
  const { data: res } = await api.get<
    SuccessResponse<HotSelectResponse>
  >(`/pdf/pdfMark/v2/web/hotSelect`, {
    params,
  });

  return res.data.annotations;
};

export const ERROR_CODE_SORT = 4000;

const heightWidthScale = Number(localStorage.getItem('hws') || 1.2);

export const convertTextAnnotation = (map: PageShapeTextMap) => {
  Object.keys(map).forEach((pageNumber) => {
    const list = map[pageNumber];
    if (!list) {
      return;
    }

    const pdfStore = usePdfStore();
    const pdfAnnotater = pdfStore.getAnnotater(currentNoteInfo.value?.noteId);
    const ipadUnit = getIpadUnit(pdfAnnotater!, Number(pageNumber) - 1);

    list.forEach((item) => {
      if (
        item.type !== IDEAAnnotateType.IDEAAnnotateTypeTextBox ||
        !item.textBox ||
        (item as any)[CONVERTED]
      ) {
        return;
      }

      item.textBox.height = (item.textBox.height * ipadUnit) / heightWidthScale;
      item.textBox.width = (item.textBox.width * ipadUnit) / heightWidthScale;
      (item as any)[CONVERTED] = true;
    });
  });
};

const restoreTextAnnotation = (item: WebNoteAnnotationModel) => {
  if (item.type !== IDEAAnnotateType.IDEAAnnotateTypeTextBox || !item.textBox) {
    return item;
  }

  const pdfStore = usePdfStore();
  const pdfAnnotater = pdfStore.getAnnotater(currentNoteInfo.value?.noteId);
  const ipadUnit = getIpadUnit(pdfAnnotater!, item.pageNumber - 1);

  const copy = {
    ...item,
    textBox: {
      ...item.textBox,
      height: (item.textBox.height * heightWidthScale) / ipadUnit,
      width: (item.textBox.width * heightWidthScale) / ipadUnit,
    },
  };

  delete (copy as any)[CONVERTED];
  return copy;
};

export const postSaveAnnotation = async (item: WebNoteAnnotationModel) => {
  const params = restoreTextAnnotation(item);

  const { data: res } = await api.post<SuccessResponse<string>>(
    `/pdf/pdfMark/v2/web/save`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );

  return res.data.uuid;
};

export const addGroupAnnotation = async (params: WebNoteAnnotationModel) => {
  const { data: res } = await api.post<SuccessResponse<string>>(
    `/pdf/pdfMark/v2/web/save`,
    params
  );

  return res.data.uuid;
};

export const deleteAnnotation = async (params: {
  annotationId: string;
  noteId: string;
  groupId: string;
}) => {
  const { data } = await api.post(
    `/pdf/pdfMark/v2/web/delete`,
    { ...params, id: params.annotationId }
  );

  return data;
};

export const editAnnotation = async (item: WebNoteAnnotationModel) => {
  const params = restoreTextAnnotation(item);

  //return api.put(`${REQUEST_SERVICE_NAME_APP}/pdfMark/v2/web/update`, params);
  return api.post(`/pdf/pdfMark/v2/web/update`, params);
};

// export const revokeDelete = async (id: string) => {
//   const { data } = await api.post<SuccessResponse<any>>(
//     `${REQUEST_SERVICE_NAME_APP}/pdfMark/revokeDelete`,
//     { id }
//   );

//   console.log(data);

//   return data;
// };

export const putSortAnnotations = async (params: SortAnnotationRequest) => {
  const { data: res } = await api.put<SuccessResponse<SortAnnotationResponse>>(
    `${REQUEST_SERVICE_NAME_APP}/pdfMark/v2/web/sortAnnotation`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );

  return res.data;
};

export const createAnnotationTag = async (params: CreateAnnotateTagRequest) => {
  const { data: res } = await api.post<SuccessResponse<CreateAnnotateTagResponse>>(
    `/pdf/marktag/save`,
    params
  );

  return res.data;
};

export const addTagToAnnotation = async (params: AddTagToAnnotateRequest) => {
  const { data: res } = await api.post(
    `/pdf/marktag/relation/mark/save`,
    params
  );

  return res;
};

export const deleteTagFromAnnotation = async (
  params: DeleteTagToAnnotateRequest
) => {
  const { data: res } = await api.post(
    `/pdf/marktag/relation/mark/delete`,
    params
  );

  return res;
};

export const paramsSerializer = (
  p: Record<string, number | string | (number | string)[]>
) => {
  const search = new URLSearchParams();

  Object.keys(p).forEach((key) => {
    const value = p[key];
    if (value instanceof Array) {
      value.forEach((item) => {
        search.append(key, String(item));
      });
    } else {
      search.append(key, String(value));
    }
  });

  return String(search);
};
