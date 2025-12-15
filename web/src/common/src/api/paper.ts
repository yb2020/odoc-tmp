import {
  CancelCollectPaperRequest,
  CollectPaperRequest,
  MyCollectedDocInfo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/userCenter/UserDoc';
import { REQUEST_SERVICE_NAME_APP } from './const';
import { $AuthorInfo } from './profile';
import { SuccessResponse } from './type';
import { AxiosRequestConfig } from 'axios';
import api from './axios';

export type $PaperAuthor = {
  id?: string;
  name: string;
};

export type $CourseItem = {
  name: string;
  id: string;
};

export type $CollectInfo = {
  id: string;
  createDate: number;
  modifyDate: number;
};

export type $AbstractInfo = {
  bestAbsAuthor: $AuthorInfo;
  bestAbsContent: string;
  bestAbsModifyTime: number;
  totalCount: number;
  noteId: string;
};

export enum RecommendReason {
  ref = 'ref',
  classic = 'classic',
  new = 'new',
  hot = 'hot',
}

export type $NoteInfo = {
  author: $AuthorInfo;
  noteCount: number;
  likeCount: number;
  modifyTime: number;
  content: string;
  noteId: string;
};

export type $PaperDetail = {
  id: string;
  authorList: $PaperAuthor[];
  title: string;
  summary: string;
  // ext: $PaperExt
  date: string;
  noteCount: number;
  courseList: $CourseItem[];
  pdfCount: number;
  pdfOwner: string;
  pdfId: any;
  pdfSourceUrl: string;
  abstractInfo: $AbstractInfo;
  isCollected: boolean;
  collectionCount: number;
  // classQuestions: {
  //   bestCqAuthor: $AuthorInfo
  //   noteId: string
  //   questionsInfo: $QuestionInfo[]
  // }
  hotNote?: {
    totalNoteCount: number;
    noteList: $NoteInfo[];
  };
  conference?: number;
  videoLink: string;
  isCollect: $CollectInfo | null;
  recentlyCollectionCount: number;
  conferenceInfo?: string;
  journal?: string;
  publishDate: number;
  venues?: string[];
  primaryVenue?: string;
  versionCount?: number;
  recommendId?: string;
  recommendItem?: {
    id?: string;
    qid?: string;
    qtitle?: string;
    reason?: RecommendReason;
  };
  venueTags?: string[];
  citationCount?: number;
};

export interface classifyItem {
  classifyId: string;
  classifyName: string;
  isContain: boolean;
}

// 收藏论文
export const collectPaper = async (
  param: Partial<CollectPaperRequest>,
  config?: AxiosRequestConfig
) => {
  const res = await api.post<SuccessResponse<MyCollectedDocInfo>>(
    `${REQUEST_SERVICE_NAME_APP}/userDoc/collectPaper`,
    param,
    config
  );
  return res.data.data || {};
};
// 取消收藏论文
export const cancelCollectPaper = async (param: CancelCollectPaperRequest) => {
  const res = await api.post<SuccessResponse<classifyItem>>(
    `${REQUEST_SERVICE_NAME_APP}/userDoc/cancelCollectPaper`,
    param
  );
  return res.data.data || {};
};
