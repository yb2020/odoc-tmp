import { defineStore } from 'pinia';
import {
  NeedVipException,
  PayStatus,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import {
  VipRole,
  VipPrivilege,
  VipType,
  VipProfileResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/user/VipUserInterface';
import { getVipProfile /* postOcrRemainCount */ } from '@common/api/vip';
import { getVipPayConfig } from '@common/api/vipPay';
import { useI18n } from 'vue-i18n';
import { getHigherVipType } from '../components/Premium/types';
import { GroupProductItem } from '../components/Premium/gropbuy'; // 引入 gropbuy

export { VipType };

export type VipTypePayable = Exclude<
  VipType,
  VipType.FREE | VipType.UNRECOGNIZED
>;

export interface ReportParams {
  page_type?: string;
  type_parameter?: string;
  element_name?: string;
  element_parameter?: string;
}

export interface LimitDialogProps {
  leftBtn?: {
    text: string;
    url: string;
  };
  exception?: NeedVipException;
  reportParams?: ReportParams;
}

export interface PayDialogProps {
  needVipType?: VipTypePayable;
  reportParams?: ReportParams;
  isGroupBuy?: boolean; // 添加是否团购参数
  groupProductItem?: GroupProductItem[];
  onPaySucc?: (x: PayStatus) => void;
  onPayCancel?: () => void;
}

export interface VipState {
  env: string;
  inited: boolean;
  enabled: boolean;
  payByDialog: boolean;
  roles: VipRole[];
  privileges: Partial<VipPrivilege>;
  showLimitDialog: '' | 'vip' | 'ocr';
  limitDialogMessage: string;
  limitDialogProps: LimitDialogProps;
  showPayDialog: boolean;
  payStatus: PayStatus;
  payDialogEnabled: boolean;
  payDialogProps: PayDialogProps;
  ocrRemainCount: number;
  profile?: VipProfileResponse;
  payOrigin?: string;
}

export const PremiumVipPreferences = {
  [VipType.FREE]: {
    key: 'free',
    icon: null,
    color: '#A8AFBA',
  },
  [VipType.STANDARD]: {
    key: 'standard',
    // icon: IconPremium,
    color: '#6AA9EC',
  },
  [VipType.PROFESSIONAL]: {
    key: 'profession',
    // icon: IconPremium,
    color: '#5A77EB',
  },
  [VipType.ENTERPRISE]: {
    key: 'enterprise',
    // icon: IconPremium,
    color: '#24304B',
  },
  [VipType.OUTSTANDING]: {
    key: 'outstanding',
    // icon: IconPremium,
    color: '#24304B',
  },
  [VipType.UNRECOGNIZED]: {} as any,
};

export const useVipStore = defineStore('vip', {
  state: (): VipState => ({
    env: '',
    inited: false,
    enabled: true,
    payByDialog: false,
    roles: [],
    profile: undefined,
    privileges: {},
    showLimitDialog: '',
    limitDialogMessage: '',
    limitDialogProps: {},
    showPayDialog: false,
    payStatus: PayStatus.PAY_PRE,
    payDialogEnabled: false,
    payDialogProps: {},
    ocrRemainCount: 0,
    payOrigin: '',
  }),
  getters: {
    role(state) {
      return (
        state.roles[0] || {
          vipType: VipType.FREE,
        }
      );
    },
    seniorTxt(state) {
      if (state.enabled) {
        const i18n = useI18n();
        return i18n.t('common.premium.senior');
      }

      return 'VIP';
    },
    seniorRole(state) {
      const role = state.roles?.[0]?.vipType ?? VipType.FREE;

      if (state.enabled) {
        return getHigherVipType(role) as VipTypePayable;
      }

      return null;
    },
    seniorRoles(state: VipState) {
      let days = 0;
      return (
        state.roles
          ?.map((role, i) => {
            const leftDays = role.vipLeftDays - days;
            days = Math.max(role.vipLeftDays, days);

            return {
              ...role,
              leftDays,
              seniorIdx: i - 1,
            };
          })
          .filter((x) => x.vipType !== VipType.FREE && x.leftDays > 0) ?? []
      );
    },
  },
  actions: {
    // 已弃用：/pay/public/getVipPrivilege接口已不再使用
    async fetchVipConfig() {
      // 使用模拟数据替代API调用
      const res = await getVipPayConfig();
      
      this.inited = true;
      this.env = res.env;
      this.enabled = !!res?.enabled;
      this.payByDialog = res?.payByDialog ?? false;
      this.payOrigin = res?.payOrigin ?? '';
    },
    // 已弃用：/vip/profile接口已不再使用
    async fetchVipProfile() {
      if (Object.keys(this.privileges).length) {
        return;
      }

      // 使用模拟数据替代API调用
      const res = await getVipProfile(); // 使用模拟实现
      this.profile = res;
      this.roles = res.vipRoles ?? [];
      this.privileges = res.vipRoles?.[0]?.vipPrivilege || {
        vipType: VipType.FREE,
      };
    },
    // 已弃用：/ts/ocr/remainCount接口已不再使用
    async fetchOcrCount() {
      // 注释掉原始实现
      // const req = postOcrRemainCount({});
      // this.ocrRemainCount = await req;
      
      console.log('vip.ts: fetchOcrCount called, but API is deprecated');
      // 设置默认值
      this.ocrRemainCount = 0;
    },
    showVipLimitDialog(message: string, limitDialogProps?: LimitDialogProps) {
      this.showLimitDialog = 'vip';
      this.setDialogParams(message, limitDialogProps);
    },
    showOcrLimitDialog(message: string, limitDialogProps?: LimitDialogProps) {
      this.setDialogParams(message, limitDialogProps);
      this.showLimitDialog = 'ocr';
    },
    setDialogParams(message: string, limitDialogProps?: LimitDialogProps) {
      this.limitDialogMessage = message;
      this.limitDialogProps = limitDialogProps || {
        exception: {
          needVipType: VipType.STANDARD,
        },
      };
    },
    hideVipLimitDialog() {
      this.showLimitDialog = '';
    },
    showVipPayDialog(props: PayDialogProps) {
      this.showPayDialog = true;
      this.showLimitDialog = '';
      this.payStatus = PayStatus.PAY_PRE;
      this.payDialogProps = props;
    },
    hideVipPayDialog(status?: PayStatus) {
      this.showPayDialog = false;
      this.payStatus = status ?? PayStatus.UNRECOGNIZED;
      this.payDialogProps = {};
    },
  },
});
