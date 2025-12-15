import {
  VipPrivilege,
  VipType,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/user/VipUserInterface';

// 定义产品item的接口
export interface GroupProductItem {
  name: string;
  purchaseTips: string;
  originalPrice: number;
  price: number;
  count: number;
  vipType: VipType;
  selected: boolean; // 添加选择状态属性
  isDiscount: boolean; // 是否参与折扣
  discount: number; // 折扣率
}
