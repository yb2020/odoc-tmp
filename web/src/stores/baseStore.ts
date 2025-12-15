import { defineStore } from 'pinia';
import { getRedDotConfig } from '~/src/api/user';
import { ResponseError } from '../api/type';

enum FunctionTypeEnum {
  COLOR_PATTERN = 9,
}

interface RedDotItem {
  functionType: number;
  [key: string]: any;
}

interface GetUserRedDotListResponse {
  redDotList?: RedDotItem[];
  [key: string]: any;
}

export interface BaseState {
  pageError: ResponseError | null;
  colorRedDot: boolean;
}

export const useBaseStore = defineStore('base', {
  state: (): BaseState => ({
    pageError: null,
    colorRedDot: false,
  }),
  actions: {
    setColorRedDot(flag: boolean) {
      this.colorRedDot = flag;
    },
    async getRedDotInfo() {
      const data = await getRedDotConfig() as GetUserRedDotListResponse;

      const flag = data?.redDotList?.some(
        (item) => item.functionType === FunctionTypeEnum.COLOR_PATTERN
      );

      if (flag) {
        this.colorRedDot = flag;
      }
    },
  },
});
