import { UserInfo } from "~/src/api/user";
import { Nullable } from "~/src/typings/global";

export interface UserState {
  userInfo: Nullable<UserInfo>;
  token?: string;
  isAuthenticated: boolean;
  showLoginDialog: boolean;
  isFetchGetTrafficApi: boolean,
  isSimpleLogin: boolean,
  lastUserInfoFetch: number, // 添加字段记录上次获取用户信息的时间戳
}