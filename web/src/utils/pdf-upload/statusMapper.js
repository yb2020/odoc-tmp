import { UserDocParsedStatusEnum } from 'go-sea-proto/gen/ts/doc/UserDocParsedStatus';

/**
 * 状态到 i18n key 的映射
 * @param {UserDocParsedStatusEnum} status - 状态枚举值
 * @returns {string} - i18n key
 */
export const getStatusI18nKey = (status) => {
  const statusKeyMap = {
    // 其他状态 (0-5)
    [UserDocParsedStatusEnum.READY]: 'home.upload.result.status.ready',
    [UserDocParsedStatusEnum.REPARSE]: 'home.upload.result.status.reparse',
    
    // 下载状态 (5-10)
    [UserDocParsedStatusEnum.DOWNLOADING]: 'home.upload.result.status.downloading',
    [UserDocParsedStatusEnum.DOWNLOADED]: 'home.upload.result.status.downloaded',
    [UserDocParsedStatusEnum.DOWNLOAD_FAILED]: 'home.upload.result.status.downloadFailed',
    
    // 上传状态 (11-20)
    [UserDocParsedStatusEnum.UPLOADING]: 'home.upload.result.status.uploading',
    [UserDocParsedStatusEnum.UPLOADED]: 'home.upload.result.status.uploaded',
    [UserDocParsedStatusEnum.UPLOAD_FAILED]: 'home.upload.result.status.uploadFailed',
    
    // 解析状态 (21-40)
    [UserDocParsedStatusEnum.GENERATING_BASE_DATA]: 'home.upload.result.status.generatingBaseData',
    [UserDocParsedStatusEnum.BASE_DATA_GENERATED]: 'home.upload.result.status.baseDataGenerated',
    [UserDocParsedStatusEnum.BASE_DATA_GENERATE_FAILED]: 'home.upload.result.status.baseDataGenerateFailed',
    [UserDocParsedStatusEnum.PARSING_HEADER_DATA]: 'home.upload.result.status.parsingHeaderData',
    [UserDocParsedStatusEnum.HEADER_DATA_PARSED]: 'home.upload.result.status.headerDataParsed',
    [UserDocParsedStatusEnum.HEADER_DATA_PARSE_FAILED]: 'home.upload.result.status.headerDataParseFailed',
    [UserDocParsedStatusEnum.PARSING_CONTENT_DATA]: 'home.upload.result.status.parsingContentData',
    [UserDocParsedStatusEnum.CONTENT_DATA_PARSED]: 'home.upload.result.status.contentDataParsed',
    [UserDocParsedStatusEnum.CONTENT_DATA_PARSE_FAILED]: 'home.upload.result.status.contentDataParseFailed',
    [UserDocParsedStatusEnum.PARSE_FAILED]: 'home.upload.result.status.parseFailed',
  };
  
  return statusKeyMap[status] || 'home.upload.result.status.unknown';
};




/**
 * 获取状态显示文本
 * @param {UserDocParsedStatusEnum} status - 状态枚举值
 * @param {Function} t - 翻译函数
 * @returns {string} - 状态显示文本
 */
export const getStatusText = (status, t) => {
  if (!t || typeof t !== 'function') {
    return `状态${status}`;
  }
  
  const key = getStatusI18nKey(status);
  
  try {
    const result = t(key);
    return result;
  } catch (error) {
    console.error("翻译过程出错", { key, status, error });
    return key;
  }
};

/**
 * 根据状态计算进度百分比
 * @param {UserDocParsedStatusEnum} status - 状态枚举值
 * @returns {number} - 进度百分比 (0-100)
 */
export const getStatusProgress = (status) => {
  const progressMap = {
    // 其他状态
    [UserDocParsedStatusEnum.READY]: 0,
    [UserDocParsedStatusEnum.REPARSE]: 0,
    
    // 下载状态
    [UserDocParsedStatusEnum.DOWNLOADING]: 10,
    [UserDocParsedStatusEnum.DOWNLOADED]: 15,
    [UserDocParsedStatusEnum.DOWNLOAD_FAILED]: 100,
    
    // 上传状态
    [UserDocParsedStatusEnum.UPLOADING]: 25,
    [UserDocParsedStatusEnum.UPLOADED]: 35,
    [UserDocParsedStatusEnum.UPLOAD_FAILED]: 100,
    
    // 解析状态
    [UserDocParsedStatusEnum.GENERATING_BASE_DATA]: 40,
    [UserDocParsedStatusEnum.BASE_DATA_GENERATED]: 50,
    [UserDocParsedStatusEnum.BASE_DATA_GENERATE_FAILED]: 100,
    [UserDocParsedStatusEnum.PARSING_HEADER_DATA]: 60,
    [UserDocParsedStatusEnum.HEADER_DATA_PARSED]: 100, // 头部数据解析完成即视为完成
    [UserDocParsedStatusEnum.PARSING_CONTENT_DATA]: 100, // 内容解析中，已经可以阅读
    [UserDocParsedStatusEnum.CONTENT_DATA_PARSED]: 100, // 内容解析完成

    
  };
  
  return progressMap[status] || 0;
};

/**
 * 判断状态是否为错误状态
 * @param {UserDocParsedStatusEnum} status - 状态枚举值
 * @returns {boolean} - 是否为错误状态
 */
export const isErrorStatus = (status) => {
  const errorStatuses = [
    UserDocParsedStatusEnum.DOWNLOAD_FAILED,
    UserDocParsedStatusEnum.UPLOAD_FAILED,
    UserDocParsedStatusEnum.BASE_DATA_GENERATE_FAILED,
    UserDocParsedStatusEnum.HEADER_DATA_PARSE_FAILED,
    UserDocParsedStatusEnum.CONTENT_DATA_PARSE_FAILED,
    UserDocParsedStatusEnum.PARSE_FAILED,
  ];
  
  return errorStatuses.includes(status);
};

/**
 * 判断状态是否为成功状态（包括头部数据解析完成）
 * @param {UserDocParsedStatusEnum} status - 状态枚举值
 * @returns {boolean} - 是否为成功状态
 */
export const isSuccessStatus = (status) => {
  return status >= UserDocParsedStatusEnum.HEADER_DATA_PARSED;
};

/**
 * 判断状态是否为进行中状态
 * @param {UserDocParsedStatusEnum} status - 状态枚举值
 * @returns {boolean} - 是否为进行中状态
 */
export const isProcessingStatus = (status) => {
  // 首先排除错误状态
  if (isErrorStatus(status)) {
    return false;
  }
  
  // 然后判断是否小于完成状态
  return status < UserDocParsedStatusEnum.HEADER_DATA_PARSED;
};

/**
 * 获取状态对应的CSS类名
 * @param {UserDocParsedStatusEnum} status - 状态枚举值
 * @returns {string} - CSS类名
 */
export const getStatusClassName = (status) => {
  if (isErrorStatus(status)) {
    return 'error';
  }
  
  if (isSuccessStatus(status)) {
    return 'finish';
  }
  
  if (isProcessingStatus(status)) {
    // 根据具体状态返回对应的CSS类名
    switch (status) {
      case UserDocParsedStatusEnum.UPLOADING:
        return 'upload';
      case UserDocParsedStatusEnum.PARSING_HEADER_DATA:
        return 'parse';
      default:
        return 'upload'; // 默认使用upload样式
    }
  }
  
  return 'waiting';
};


/**
 * 判断当前状态是否是解析失败状态
 * @param {UserDocParsedStatusEnum} status - 状态枚举值
 * @returns {boolean} - 是否为解析失败状态
 */
export const isParseFailedStatus = (status) => {
  return status === UserDocParsedStatusEnum.PARSE_FAILED || 
    status === UserDocParsedStatusEnum.CONTENT_DATA_PARSE_FAILED || 
    status === UserDocParsedStatusEnum.HEADER_DATA_PARSE_FAILED || 
    status === UserDocParsedStatusEnum.BASE_DATA_GENERATE_FAILED;
};

/**
 * 根据 parsedStatus 和 embeddingStatus 计算进度百分比
 * @param {number} parsedStatus - 解析状态枚举值
 * @param {number} embeddingStatus - embedding 状态枚举值
 * @returns {string} - 进度百分比字符串，如 '60%'
 */
export const calculateParsedProgress = (parsedStatus, embeddingStatus) => {
  let progress = '0%';

  switch (parsedStatus) {
    case UserDocParsedStatusEnum.READY:
      progress = '0%';
      break;
    case UserDocParsedStatusEnum.REPARSE:
      progress = '0%';
      break;
    case UserDocParsedStatusEnum.DOWNLOADING:
      progress = '5%';
      break;
    case UserDocParsedStatusEnum.DOWNLOADED:
      progress = '10%';
      break;
    case UserDocParsedStatusEnum.DOWNLOAD_FAILED:
      progress = '15%';
      break;
    case UserDocParsedStatusEnum.UPLOADING:
      progress = '20%';
      break;
    case UserDocParsedStatusEnum.UPLOADED:
      progress = '25%';
      break;
    case UserDocParsedStatusEnum.UPLOAD_FAILED:
      progress = '30%';
      break;
    case UserDocParsedStatusEnum.GENERATING_BASE_DATA:
      progress = '40%';
      break;
    case UserDocParsedStatusEnum.BASE_DATA_GENERATED:
      progress = '45%';
      break;
    case UserDocParsedStatusEnum.BASE_DATA_GENERATE_FAILED:
      progress = '40%';
      break;
    case UserDocParsedStatusEnum.PARSING_HEADER_DATA:
      progress = '50%';
      break;
    case UserDocParsedStatusEnum.HEADER_DATA_PARSED:
      progress = '60%';
      break;
    case UserDocParsedStatusEnum.HEADER_DATA_PARSE_FAILED:
      progress = '50%';
      break;
    case UserDocParsedStatusEnum.PARSING_CONTENT_DATA:
      progress = '70%';
      break;
    case UserDocParsedStatusEnum.CONTENT_DATA_PARSED:
      progress = '80%';
      // 如果 embedding 已完成，则进度为 100%
      if (embeddingStatus === UserDocParsedStatusEnum.EMBEDDED) {
        progress = '100%';
      }
      break;
    case UserDocParsedStatusEnum.CONTENT_DATA_PARSE_FAILED:
      progress = '70%';
      break;
    case UserDocParsedStatusEnum.PARSE_FAILED:
      progress = '0%';
      break;
    default:
      // unknown 或其他未知状态
      progress = '0%';
      break;
  }

  return progress;
};
