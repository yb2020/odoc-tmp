import { api } from '@common/api/axios';
import { AuthorizationLoginRequest } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/user/UserAuthorization';
import { SuccessResponse } from '@common/api/type';
import { REQUEST_APPID, REQUEST_SERVICE_NAME_USER } from '@common/api/const';

export const getAuthorizationCode = async () => {
  const res = await api.get<SuccessResponse<string>>(
    `${REQUEST_SERVICE_NAME_USER}/user/authorization/authorizationCode`
  );
  return res.data.data;
};

export const getAuthorizationUrl = (
  p: Omit<AuthorizationLoginRequest, 'userId'> & {
    env?: string;
  },
  origin = window.location.origin
) => {
  const url = new URL(
    `${origin}/api${REQUEST_SERVICE_NAME_USER}/oauth/authorization/logon`
  );
  url.searchParams.set('appId', REQUEST_APPID);
  Object.keys(p).forEach((key) => {
    const v = p[key as keyof typeof p];
    if (v) {
      url.searchParams.set(key, v);
    }
  });
  return url.toString();
};
