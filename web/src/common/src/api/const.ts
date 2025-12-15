export const REQUEST_SERVICE_NAME_SEARCH = '/microService-aiKnowledge-search';
export const REQUEST_SERVICE_NAME_APP = '/microService-app-aiKnowledge';
export const REQUEST_SERVICE_NAME_USER = '/microService-user';
export const REQUEST_SERVICE_NAME_ACL = '/microService-acl';
export const REQUEST_SERVICE_NAME_PAY = '/readpaper-pay';
export const REQUEST_SERVICE_NAME_SUPERVISOR_REVIEW =
  '/microservice-supervisor-review';
export const REQUEST_SERVICE_READPAPER_DOC = '/microservice-readpaper-doc';
export const REQUEST_SERVICE_NAME_TRANSLATE = '/microService-translate';
export const REQUEST_SERVICE_NAME_AI_POLISH = '/microService-aiPolish';
export const REQUEST_SERVICE_NAME_AI_REVIEW = '/microservice-ai-review';
export const REQUEST_SERVICE_NAME_AI_READING = '/microservice-ai-reading';

export const UNKNOWN_ERROR_CODE = 999;
export const UNKNOWN_ERROR_MESSAGE = 'unknown error';

// 50008: Illegal token; 50012: Other clients logged in; 50014: Token expired;

export const ERROR_CODE_UNLOGIN = 1000;
export const ERROR_CODE_ILLEGAL_TOKEN = 50008;
export const ERROR_CODE_REPEAT_LOGIN = 50012;
export const ERROR_CODE_TOKEN_EXPIRED = 50014;
export const SUCCESS_CODE = 0;
export const ERROR_CODE_NEED_VIP = 2;
export const ERROR_CODE_BEANS_NOT_ENOUGH = 3;
export const ERROR_CODE_BEANS_CASH_DEDUCTION = 4;

export const HEADER_CANCLE_AUTO_ERROR = 'x-custom-handle-error';
export const getHeadersWithCancelAutoError = () => ({
  headers: {
    [HEADER_CANCLE_AUTO_ERROR]: true,
  },
});

export const ERROR_CODE_DOC_LIMIT = 2;

export const REQUEST_APPID = 'aiKnowledge';

export const REQUEST_ORGID = '535992339879038976';

// 非法的身份刷新(使用refreshToken刷新accessToken时的错误，refreshToken验证失败，弹登录框)
export const ERROR_CODE_ILLEGAL_IDENTITY_REQUEST_REFRESH = 50010;

// 过期的身份请求刷新token(使用refreshToken刷新accessToken时的错误，表示refreshToken已过期，弹登录框)
export const ERROR_CODE_EXPIRED_IDENTITY_REQUEST_REFRESH = 50011;

export const COOKIE_REFRESH_TOKEN = 'aiKnowledge-refresh-token';
