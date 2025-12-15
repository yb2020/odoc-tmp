export class ResponseError extends Error {
  code: number;
  extra?: object;

  constructor({
    code,
    message,
    extra,
  }: {
    code: number;
    message: string;
    extra?: object;
  }) {
    super(message);
    this.code = code;
    this.extra = extra;
  }
}
export interface SuccessResponse<T> {
  status: 1;
  data: T;
  message: string;
}

export interface FailResponse<T = unknown> {
  code?: number;
  status: number;
  message: string;
  data: T;
}

export type Response<T> = SuccessResponse<T> | FailResponse<T>;

export function isSuccessResponse<T>(
  res: SuccessResponse<T> | FailResponse
): res is SuccessResponse<T> {
  return (res as SuccessResponse<T>).status === 1;
}

export type RequestParam<T> = Omit<T, 'pdfId' | 'paperId' | 'noteId'>;

export interface RefreshTokenReq {
  refreshToken: string;
  scope: string;
}

export interface RefreshTokenResp {
  accessToken: string;
  // accessToken有效期
  expiresAt: number;
  // refreshToken
  refreshToken: string;
}
