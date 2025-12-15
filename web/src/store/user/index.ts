import {
  Module,
  Store as VuexStore,
  GetterTree,
  ActionTree,
  ActionContext,
  MutationTree,
} from 'vuex';
import { User } from 'go-sea-proto/gen/ts/user/User'
import { getUserInfo } from '../../api/user';
import { RootState } from '../types';
import { UserState } from './type';
// import { GetTrafficResp } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/ab/ABTestInfo'
import reporter from '@idea/aiknowledge-report';

type Getters = {
  isLogin(state: UserState): boolean;
};

const getters: GetterTree<UserState, RootState> & Getters = {
  isLogin: (state) => state.isAuthenticated
};

export enum UserMutationTypes {
  SET_USERINFO = 'setUserInfo',
  SHOW_LOGIN_DIALOG = 'showLoginDialog',
  UPDATE_LAST_USER_INFO_FETCH = 'updateLastUserInfoFetch',
}

export type Mutations<S = UserState> = {
  [UserMutationTypes.SET_USERINFO](state: S, payload: User | null): void;
  [UserMutationTypes.SHOW_LOGIN_DIALOG](state: S, payload: {
    show: boolean
    trafficData?: any | null | undefined
  }): void;
  [UserMutationTypes.UPDATE_LAST_USER_INFO_FETCH](state: S, payload: number): void;
};

const mutations: MutationTree<UserState> & Mutations = {
  [UserMutationTypes.SET_USERINFO](state, payload) {
    if (payload) {
      state.userInfo = payload;
      state.isAuthenticated = true;
    } else {
      state.userInfo = null;
      state.isAuthenticated = false;
    }
  },
  [UserMutationTypes.SHOW_LOGIN_DIALOG](state, payload) {
    state.isFetchGetTrafficApi = true

    if (payload?.trafficData?.item?.strategyId === 'login_page_simple') {
      state.isSimpleLogin = true
    }
 

    state.showLoginDialog = payload.show;
  },
  [UserMutationTypes.UPDATE_LAST_USER_INFO_FETCH](state, payload) {
    state.lastUserInfoFetch = payload;
  },
};

type AugmentedActionContext = {
  commit<K extends keyof Mutations>(
    key: K,
    payload: Parameters<Mutations[K]>[1]
  ): ReturnType<Mutations[K]>;
} & Omit<ActionContext<UserState, RootState>, 'commit'>;

export enum UserActionTypes {
  GET_USERINFO = 'getUserInfo',
  GO_LOGIN = 'goLogin'
}

interface Actions {
  [UserActionTypes.GET_USERINFO]: (
    context: AugmentedActionContext
  ) => Promise<void>;
  [UserActionTypes.GO_LOGIN]: (
    context: AugmentedActionContext,
    show: boolean,
  ) => Promise<void>;
}

const actions: ActionTree<UserState, RootState> & Actions = {
  async [UserActionTypes.GET_USERINFO]({ commit, state }) {
    const now = Date.now();
    const DEBOUNCE_TIME = 5000; // 5秒内不重复请求
    
    // 如果上次请求时间在 5 秒内，并且用户已经认证，则不再重复请求
    if (
      state.lastUserInfoFetch > 0 && 
      now - state.lastUserInfoFetch < DEBOUNCE_TIME && 
      state.isAuthenticated
    ) {
      console.log('用户信息请求已防抖，跳过重复请求');
      return;
    }
    
    try {
      // 更新最后请求时间
      commit(UserMutationTypes.UPDATE_LAST_USER_INFO_FETCH, now);
      
      const res = await getUserInfo();
      if (res) {
        commit(UserMutationTypes.SET_USERINFO, res);
      } else {
        // 如果没有用户信息，确保设置为未登录状态
        commit(UserMutationTypes.SET_USERINFO, null);
      }
    } catch (error) {
      console.error('获取用户信息失败:', error);
      commit(UserMutationTypes.SET_USERINFO, null);
    }
  },
  async [UserActionTypes.GO_LOGIN]({ commit, state }, show ) {
    if (show && !state.isFetchGetTrafficApi) {
      // 注释掉AB测试相关代码，避免404错误
      // const data = await getTraffic({ experimentId: '' })
      const data = null;

      // if(data?.item) {
      //   reporter.abConfigItem = Promise.resolve(data?.item)
      // }

      commit(UserMutationTypes.SHOW_LOGIN_DIALOG, {
        show,
        trafficData: data,
      })

      return
    }

    commit(UserMutationTypes.SHOW_LOGIN_DIALOG, {
      show,
    })
  },
};

export type UserStore<S = UserState> = Omit<VuexStore<S>, 'getters'> & {
  getters: {
    [K in keyof Getters]: ReturnType<Getters[K]>;
  };
};

export const UserModule: Module<UserState, RootState> = {
  namespaced: true,
  state: () => ({
    userInfo: null,
    isAuthenticated: false,
    showLoginDialog: false,
    isFetchGetTrafficApi: false,
    isSimpleLogin: false,
    lastUserInfoFetch: 0,
  }),
  getters,
  actions,
  mutations,
};
