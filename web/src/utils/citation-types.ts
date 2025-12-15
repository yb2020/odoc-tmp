import Cite from 'citation-js'

// 基于原项目的复杂类型定义
export type CSLDateParts = number[];
export type CSLDate = Array<{
  'date-parts': CSLDateParts[];
}> | {
  'date-parts': CSLDateParts[];
} | {
  literal?: string;
  raw?: string;
};

// 简化的作者信息类型
export interface AuthorInfo {
  id?: string;
  given?: string;
  family?: string;
  literal: string;
  isAuthentication?: boolean;
}

// 基于原项目的PaperData类型（简化版）
export type PaperData = {
  [k: string]: undefined | any;
  pdfId?: string;
  paperId?: string;
  docName: string;
  authors?: string[];
  authorList?: AuthorInfo[];
  venues?: string[];
  primaryVenue?: string;
  publishDate?: string;
  language?: string;
  page?: string;
  issue?: string;
  volume?: string;
  eventTitle?: string;
  eventPlace?: string;
  doi?: string;
  url?: string;
  // 显示相关字段（简化）
  displayAuthor?: { authorInfos: AuthorInfo[]; userEdited: boolean };
  displayVenue?: { venueInfos: string[]; userEdited: boolean };
  displayPublishDate?: { publishDate: string; userEdited: boolean };
};

// OutputFormat 基于 citeproc Engine 的输出格式
export type OutputFormat = 'html' | 'text' | 'rtf';

// RENDER_ERR 枚举
export enum RENDER_ERR {
  NO_BIBL_STYLE = 1,
}

// 引文渲染状态
export enum CitationRenderStatus {
  failed = 0,
  successful = 1,
  empty = 2,
}

// 本地实现 paper2MetaData
export function paper2MetaData(paperData: PaperData) {
  const authors = paperData.authorList || paperData.authors?.map(name => ({ literal: name })) || [];
  const displayAuthors = paperData.displayAuthor?.authorInfos || authors;
  
  const venue = paperData.displayVenue?.venueInfos?.[0] || paperData.primaryVenue || paperData.venues?.[0] || '';
  const publishDate = paperData.displayPublishDate?.publishDate || paperData.publishDate;
  
  let issuedDate: CSLDate | undefined;
  if (publishDate) {
    const date = new Date(publishDate);
    if (!isNaN(date.getTime())) {
      issuedDate = { 'date-parts': [[date.getFullYear(), date.getMonth() + 1, date.getDate()]] };
    }
  }

  const metaData = {
    id: paperData.paperId || paperData.pdfId || `paper_${Date.now()}`,
    type: venue ? 'article-journal' : 'article',
    title: paperData.docName || '',
    author: displayAuthors.map((author: any) => {
      // 如果有完整的given和family，使用标准格式
      if (author.given && author.family) {
        return {
          given: author.given,
          family: author.family
        };
      }
      // 否则使用literal格式（避免同时设置given/family和literal）
      return {
        literal: author.literal || author
      };
    }),
    issued: issuedDate,
    'container-title': venue,
    page: paperData.page,
    volume: paperData.volume,
    issue: paperData.issue,
    DOI: paperData.doi,
    URL: paperData.url,
    language: paperData.language,
    'event-title': paperData.eventTitle,
    'event-place': paperData.eventPlace,
  };

  // 移除空值
  Object.keys(metaData).forEach(key => {
    if (metaData[key as keyof typeof metaData] === undefined || metaData[key as keyof typeof metaData] === '') {
      delete metaData[key as keyof typeof metaData];
    }
  });

  return metaData;
}

// 简化的 citeprocSys 实现
export const citeprocSys = {
  addItem: (id: string, data: any, force?: boolean) => {
    console.log('Adding citation item:', { id, title: data.title });
  }
};

// 改进的 renderBibl 实现
export async function renderBibl(
  style: string, 
  arg1: any[], 
  arg2: any[], 
  items: any[], 
  outputFormat: OutputFormat, 
  lang: string
) {
  try {
    const cite = new Cite(items);
    let result: string;
    
    // 处理不同的引文样式
    if (style.toLowerCase().includes('bibtex')) {
      result = cite.format('bibtex');
    } else if (style.toLowerCase().includes('apa')) {
      result = cite.format('bibliography', {
        format: outputFormat,
        template: 'apa',
        lang: lang
      });
    } else if (style.toLowerCase().includes('mla')) {
      result = cite.format('bibliography', {
        format: outputFormat,
        template: 'mla',
        lang: lang
      });
    } else {
      // 默认使用指定的样式
      result = cite.format('bibliography', {
        format: outputFormat,
        template: style,
        lang: lang
      });
    }
    
    return {
      txt: result,
      msgCode: null
    };
  } catch (error) {
    console.error('Citation rendering failed:', error);
    return {
      txt: '',
      msgCode: RENDER_ERR.NO_BIBL_STYLE
    };
  }
}
