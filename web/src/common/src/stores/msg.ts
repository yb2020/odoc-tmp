import { defineStore } from 'pinia';
import {
  DetailMessageCount,
  MessageTypeEnum,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/operator/admin/msg/SystemMsgInfo';
import useSharedWebsocket from '@common/hooks/useSharedWebsocket';
import { getRedDotConfig, getRedDotInfo } from '~/src/api/user';
import { watch } from 'vue';
import { FunctionTypeEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/RedDot';

export interface MsgState {
  // disabled: boolean
  msgCount: number;
  subMsgCountMap: {
    [k in MessageTypeEnum]: number;
  };
  redDotMap: {
    [k in FunctionTypeEnum]: boolean;
  };
}

export const useMsgStore = defineStore('msg', {
  state: () =>
    ({
      // disabled: false,
      msgCount: 0,
      subMsgCountMap: {},
      redDotMap: {},
    }) as MsgState,
  actions: {
    async init() {
      const data = await getRedDotInfo();
      if (data) {
        this.updateMsgCount(data.msgCount);
        this.updateSubMsgCount(data.detailMessageCounts);
      }
    },
    async initRedDot() {
      const data = await getRedDotConfig();

      data?.redDotList?.forEach((item) => {
        this.redDotMap[item.functionType] = true;
      });
    },
    initWS() {
      const { msg } = useSharedWebsocket();

      watch(msg, () => {
        if (msg.value?.scene === 'message_center') {
          this.updateMsgCount(msg.value.content);
        } else if (msg.value?.scene === 'OFFICIAL_MESSAGE') {
          this.updateSubMsgCount(msg.value.detailMessageCounts!);
        }
      });
    },
    updateMsgCount(count: number) {
      this.msgCount = count;
    },
    updateSubMsgCount(counts: DetailMessageCount[]) {
      if (Array.isArray(counts)) {
        this.subMsgCountMap = counts.reduce(
          (map, { messageTypeEnum: type, count }) => {
            map[type] = count;
            return map;
          },
          {} as { [k in MessageTypeEnum]: number }
        );
      }
    },
  },
});
