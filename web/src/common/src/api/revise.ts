import {
  ProjectRequest,
  ProjectResponse,
  PageRequest,
  PageResponse,
  PageInfoRequest,
  PageInfoResponse,
  AllPolishTaskRequest,
  AllPolishTaskResponse,
  PolishResultRequest,
  PolishAllResult,
  ReviewResult,
  PolishTitleResult,
  PolishAbstractResult,
  PolishIntroductionResult,
  PolishRequest,
  PolishResponse,
  PolishTaskType,
  SavePageInfoRequest,
  SavePageInfoResponse,
  RepolishTaskRequest,
  RepolishTaskResponse,
  PolishAcceptRequest,
  DeleteTaskRequest,
  PolishIgnoreRequest,
  CreateProjectResponse,
  RecentProjectsResponse,
  DeleteProjectRequest,
  DeleteProjectResponse,
  GetSegmentPreferencesResponse,
  SaveSegmentPreferenceRequest,
  UpdateProjectRequest,
  UpdateProjectResponse,
  TitleGenResult,
  AbsGenResult,
  DeleteSentenceRequest,
  AbsModifyResult,
  GetApplyQuotaReq,
  GetApplyQuotaResp,
  ApplyQuotaReq,
  ApplyQuotaResp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import {
  TagsRequest,
  TagsResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishPdf';
import {
  TrialCountResponse,
  TrialIsEnabledReq,
  TrialIsEnabledResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/polish/TrialCount';
import { api } from '@common/api/axios';
import { SuccessResponse } from '@common/api/type';
import {
  HEADER_CANCLE_AUTO_ERROR,
  REQUEST_SERVICE_NAME_AI_POLISH,
} from '@common/api/const';
import { BeanScenes } from '@common/hooks/useAIBeans';
import { AxiosRequestConfig } from 'axios';

export const ERR_CODE_RESOURCE_LIMIT = 400001;

export const ERR_CODE_RESOURCE_NETWORK_ERROR = 400002;

export enum PolishStrType {
  OVERALL = 'OVERALL',
  GEN_TITLE = 'GEN_TITLE',
  ABSTRACT = 'ABSTRACT',
  GEN_ABS = 'GEN_ABS',
  MODIFY_ABS = 'MODIFY_ABS',
  INTRODUCTION = 'INTRODUCTION',
  RELATED_WORK = 'RELATED_WORK',
}

export const PolishType2BeanScene = {
  [PolishTaskType.REVIEW]: BeanScenes.REVIEW,
  [PolishTaskType.ALL_POLISH]: BeanScenes.POLISH_ALL,
  [PolishTaskType.TITLE_POLISH]: BeanScenes.TITLE_GEN,
  [PolishTaskType.TITLE_GEN]: BeanScenes.TITLE_GEN,
  [`${PolishTaskType.TITLE_GEN}2`]: BeanScenes.TITLE_GEN,
  [PolishTaskType.ABSTRACT_POLISH]: BeanScenes.ABSTRACT_POLISH_1,
  [`${PolishTaskType.ABSTRACT_POLISH}2`]: BeanScenes.ABSTRACT_POLISH_2,
  [PolishTaskType.ABSTRACT_MODIFY]: BeanScenes.ABSTRACT_POLISH_1,
  [`${PolishTaskType.ABSTRACT_MODIFY}2`]: BeanScenes.ABSTRACT_POLISH_2,
  [PolishTaskType.ABSTRACT_GEN]: BeanScenes.ABSTRACT_GEN,
  [`${PolishTaskType.ABSTRACT_GEN}2`]: BeanScenes.ABSTRACT_GEN,
  [PolishTaskType.INTRODUCTION_POLISH]: BeanScenes.INTRO_POLISH_1,
  [`${PolishTaskType.INTRODUCTION_POLISH}2`]: BeanScenes.INTRO_POLISH_2,
};

const TASK_URL_MAP = {
  [PolishTaskType.ALL_POLISH]: `/all`,
  [PolishTaskType.REVIEW]: `/review`,
  [PolishTaskType.TITLE_POLISH]: `/title`,
  [PolishTaskType.TITLE_GEN]: `/genTitle`,
  [PolishTaskType.ABSTRACT_POLISH]: `/abstract`,
  [PolishTaskType.ABSTRACT_MODIFY]: `/absModify`,
  [PolishTaskType.ABSTRACT_GEN]: `/genAbs`,
  [PolishTaskType.INTRODUCTION_POLISH]: `/introduction`,
};

type StringifyPolishType<T> = Omit<T, 'polishType'> & {
  polishType: PolishStrType;
};

/**
 * @deprecated
 */
export const getApplyQuotaStatus = async (p?: GetApplyQuotaReq) => {
  const res = await api.get<SuccessResponse<GetApplyQuotaResp>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/getApplyQuotaStatus`,
    {
      params: p,
    }
  );
  return res.data.data;
};

/**
 * @deprecated
 */
export const applyQuota = async (p?: ApplyQuotaReq) => {
  const res = await api.post<SuccessResponse<ApplyQuotaResp>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/applyQuota`,
    p
  );
  return res.data.data;
};

export const createProject = async (
  params?: object,
  config?: AxiosRequestConfig
) => {
  const res = await api.post<SuccessResponse<CreateProjectResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/createProject`,
    params,
    config
  );
  return res.data.data || {};
};

export const updateProject = async (params?: UpdateProjectRequest) => {
  const res = await api.post<SuccessResponse<UpdateProjectResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/updateProject`,
    params
  );
  return res.data.data || {};
};

export const deleteProject = async (params: DeleteProjectRequest) => {
  const res = await api.post<SuccessResponse<DeleteProjectResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/deleteProject`,
    params
  );
  return res.data.data || {};
};

export const getRecentProjects = async (params?: object) => {
  const res = await api.get<SuccessResponse<RecentProjectsResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/recentProjects`,
    { params }
  );
  return res.data.data || {};
};

/**
 * @deprecated
 */
export const getProjects = async (params: ProjectRequest) => {
  const res = await api.get<SuccessResponse<ProjectResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/getProjectIds`,
    { params }
  );
  return res.data.data || {};
};

export const getProjectVersions = async (params: PageRequest) => {
  const res = await api.get<SuccessResponse<PageResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/getPageIds`,
    {
      params,
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );
  return res.data.data || {};
};

export const getVersionData = async (params: PageInfoRequest) => {
  const res = await api.get<SuccessResponse<PageInfoResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/getPageInfo`,
    { params }
  );
  return res.data.data || {};
};

export const saveVersionData = async (
  params: Partial<SavePageInfoRequest & { parser: boolean }>
) => {
  const res = await api.post<SuccessResponse<SavePageInfoResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/savePageInfo`,
    params,
    {
      timeout: 30000,
    }
  );
  return res.data.data || {};
};

export const getProjectTasks = async (params: AllPolishTaskRequest) => {
  const res = await api.get<SuccessResponse<AllPolishTaskResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/getTasks`,
    {
      params,
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );
  return res.data.data || {};
};

export type TaskResult =
  | PolishAllResult
  | ReviewResult
  | PolishTitleResult
  | TitleGenResult
  | AbsModifyResult[]
  | AbsGenResult
  | PolishAbstractResult[]
  | PolishIntroductionResult[];

export type ArrayTypes<T> = T extends any[] ? T : never;
export type ArrayItemTypes<T> = T extends Array<infer R> ? R : never;
export type SingleTaskResult =
  | Exclude<TaskResult, ArrayTypes<TaskResult>>
  | ArrayItemTypes<TaskResult>;

export const getTaskResult = async <T extends SingleTaskResult>(
  p: PolishResultRequest,
  type: Exclude<PolishTaskType, PolishTaskType.UNRECOGNIZED>
) => {
  const res = await api.get<SuccessResponse<T>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/result${TASK_URL_MAP[type]}`,
    {
      params: p,
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );
  return res.data.data || {};
};

export const getPolishTags = async (p: StringifyPolishType<TagsRequest>) => {
  const res = await api.get<SuccessResponse<TagsResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/tags`,
    {
      params: p,
    }
  );
  return res.data.data || {};
};

export const startPolishTask = async (
  params: Partial<PolishRequest>,
  type: Exclude<PolishTaskType, PolishTaskType.UNRECOGNIZED>,
  suffix = ''
) => {
  const res = await api.post<SuccessResponse<PolishResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish${TASK_URL_MAP[type]}${suffix}`,
    params,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );
  return res.data.data || {};
};

export const startRePolishTask = async (params: RepolishTaskRequest) => {
  const res = await api.post<SuccessResponse<RepolishTaskResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/repolish/start`,
    params
  );
  return res.data.data || {};
};

export const updateRePolishTaskResult = async (
  type:
    | PolishTaskType.ABSTRACT_GEN
    | PolishTaskType.ABSTRACT_MODIFY
    | PolishTaskType.ABSTRACT_POLISH
    | PolishTaskType.INTRODUCTION_POLISH,
  {
    currentTxt,
    polishTxt,
    ...rest
  }: {
    taskId: string;
    sectionId: string;
    currentTxt: string;
    polishTxt: string;
  }
) => {
  const k = {
    [PolishTaskType.ABSTRACT_GEN]: 'Abstract',
    [PolishTaskType.ABSTRACT_MODIFY]: 'Abstract',
    [PolishTaskType.ABSTRACT_POLISH]: 'Abstract',
    [PolishTaskType.INTRODUCTION_POLISH]: 'Introduction',
  }[type];
  const p = {
    ...rest,
    [`source${k}`]: currentTxt,
    [`new${k}`]: polishTxt,
  };
  const fn = {
    [PolishTaskType.ABSTRACT_GEN]: updateRePolishAbstractResult,
    [PolishTaskType.ABSTRACT_MODIFY]: updateRePolishAbstractResult,
    [PolishTaskType.ABSTRACT_POLISH]: updateRePolishAbstractResult,
    [PolishTaskType.INTRODUCTION_POLISH]: updateRePolishIntroductionResult,
  }[type];

  return fn(p);
};

export const updateRePolishAbstractResult = async (
  params: Partial<PolishAbstractResult>
) => {
  const res = await api.put<SuccessResponse<object>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/updateAbstract`,
    params
  );
  return res.data.data;
};

export const updateRePolishIntroductionResult = async (
  params: Partial<PolishIntroductionResult>
) => {
  const res = await api.put<SuccessResponse<object>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/updateIntroduction`,
    params
  );
  return res.data.data;
};

export const acceptPolishResult = async (p: PolishAcceptRequest) => {
  const res = await api.put<SuccessResponse<object>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/accept`,
    p
  );
  return res.data.data;
};

export const ignorePolishResult = async (p: PolishIgnoreRequest) => {
  const res = await api.put<SuccessResponse<object>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/ignore`,
    p
  );
  return res.data.data;
};

export const deletePolishResult = async (p: DeleteTaskRequest) => {
  const res = await api.post<SuccessResponse<object>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/deleteTask`,
    p
  );
  return res.data.data;
};

export const deletePolishSentence = async (p: DeleteSentenceRequest) => {
  const res = await api.post<SuccessResponse<object>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/deleteSentence`,
    p
  );
  return res.data.data;
};

export const getSegmentPreferences = async (p?: object) => {
  const res = await api.get<SuccessResponse<GetSegmentPreferencesResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/getSegmentPreference`,
    p
  );
  return res.data.data;
};

export const setSegmentPreferences = async (
  p: SaveSegmentPreferenceRequest
) => {
  const res = await api.post<SuccessResponse<object>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/saveSegmentPreference`,
    p,
    {
      headers: {
        [HEADER_CANCLE_AUTO_ERROR]: 'true',
      },
    }
  );
  return res.data.data;
};

export const getTrialCount = async () => {
  const res = await api.get<SuccessResponse<TrialCountResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/trial/count`
  );
  return res.data.data;
};

export const getTrialBeta = async (params: TrialIsEnabledReq) => {
  const res = await api.get<SuccessResponse<TrialIsEnabledResponse>>(
    `${REQUEST_SERVICE_NAME_AI_POLISH}/trial/isEnabled`,
    {
      params,
    }
  );
  return res.data.data;
};
