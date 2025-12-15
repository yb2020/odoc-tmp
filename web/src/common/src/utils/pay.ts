/* eslint-disable camelcase */
import {
  PayStatus,
  WechatJsPayInfo,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';

export function isPaying(status?: PayStatus) {
  return [PayStatus.PAY_PRE, PayStatus.PAY_WAITING].includes(status!);
}

export function isWaiting(status?: PayStatus) {
  return [PayStatus.PAY_WAITING].includes(status!);
}

export function formatPrice(price: number, decimal = 2, keepzero = false) {
  const s = (price / 100).toFixed(decimal);

  return keepzero ? s : s.replace(/\.?0*$/, '');
}

declare global {
  interface Window {
    AlipayJSBridge: any;
    WeixinJSBridge: any;
  }
}

function waitAliJsSDKReady(callback: () => void) {
  if (window.AlipayJSBridge) {
    callback && callback();
  } else {
    document.addEventListener('AlipayJSBridgeReady', callback, false);
  }
}

export function callAliPay(tradeNO: string) {
  return new Promise((resolve, reject) => {
    waitAliJsSDKReady(() => {
      // 通过传入交易号唤起快捷调用方式(注意tradeNO大小写严格)
      window.AlipayJSBridge.call(
        'tradePay',
        { tradeNO },
        function (data: { memo: string; resultCode: string }) {
          if (data.resultCode === '9000') {
            resolve(true);
          } else {
            reject(
              new Error(`${data.memo || 'unknown error'}(${data.resultCode})`)
            );
          }
        }
      );
    });
  });
}

function waitWxJsSDKReady(callback: () => void) {
  if (typeof window.WeixinJSBridge === 'undefined') {
    if (document.addEventListener) {
      document.addEventListener('WeixinJSBridgeReady', callback, false);
    }
  } else {
    callback();
  }
}

export function callWechatPay(
  params: Omit<WechatJsPayInfo, 'packages'> & {
    package: string;
  }
) {
  return new Promise((resolve, reject) => {
    waitWxJsSDKReady(() =>
      window.WeixinJSBridge.invoke(
        'getBrandWCPayRequest',
        params,
        function (res: { err_msg: string }) {
          if (res.err_msg === 'get_brand_wcpay_request:ok') {
            // 使用以上方式判断前端返回,微信团队郑重提示：
            // res.err_msg将在用户支付成功后返回ok，但并不保证它绝对可靠。
            resolve(true);
            return;
          }

          reject(new Error(res.err_msg || 'unknown error'));
        }
      )
    );
  });
}
