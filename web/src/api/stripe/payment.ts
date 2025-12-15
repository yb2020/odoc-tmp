import { loadStripe } from '@stripe/stripe-js';
import { createCheckoutSession, getStripePublishableKey } from '@/api/pay';
import { createOrder } from '@/api/order';
import { Ref } from 'vue';

export interface PaymentOptions {
  orderType: number;
  numberCount: number;
  openInNewTab?: boolean;
}

export interface PaymentResult {
  success: boolean;
  error?: string;
  sessionId?: string;
}

/**
 * 处理支付按钮点击事件
 * @param options 支付选项
 * @param isLoading 加载状态的响应式引用
 * @param errorMessage 错误信息的响应式引用
 * @returns 支付处理的Promise
 */
export const handlePaymentClick = async (
  options: PaymentOptions,
  isLoading: Ref<boolean>,
  errorMessage: Ref<string | null>
): Promise<void> => {
  if (isLoading.value) return;
  
  isLoading.value = true;
  errorMessage.value = null;

  try {
    const result = await processStripePayment(options);
    
    if (!result.success) {
      errorMessage.value = result.error || '支付处理失败';
      // 5秒后清除错误信息
      setTimeout(() => {
        errorMessage.value = null;
      }, 5000);
    }
    // 成功情况下，用户已被重定向到支付页面，不需要额外处理
  } catch (err: any) {
    console.error('支付流程遇到问题:', err);
    errorMessage.value = err.message || '处理支付请求时发生了一个错误。';
    // 5秒后清除错误信息
    setTimeout(() => {
      errorMessage.value = null;
    }, 5000);
  } finally {
    isLoading.value = false;
  }
};

/**
 * 处理Stripe支付流程
 * @param options 支付选项
 * @returns 支付结果
 */
export const processStripePayment = async (options: PaymentOptions): Promise<PaymentResult> => {
  try {
    // 1.获取 publishableKey
    const publishableKeyResponse = await getStripePublishableKey({});
    if (!publishableKeyResponse || !publishableKeyResponse.publishableKey) {
      return { success: false, error: '未能从后端获取有效的 Publishable Key。' };
    }
    const publishableKey = publishableKeyResponse.publishableKey;
    
    // 2. 创建订单
    const orderResponse = await createOrder({
      orderType: options.orderType,
      numberCount: options.numberCount,
    });

    if (!orderResponse || !orderResponse.orderId) {
      return { success: false, error: '未能从后端获取有效的订单信息。' };
    }
    
    // 3. 调用后端创建 Checkout Session
    const sessionResponse = await createCheckoutSession({
      orderId: orderResponse.orderId,
    });

    if (!sessionResponse || !sessionResponse.sessionId) {
      return { success: false, error: '未能从后端获取有效的 Session ID。' };
    }

    const sessionId = sessionResponse.sessionId;

    // 4. 获取 Stripe.js 实例
    const stripe = await loadStripe(publishableKey);
    if (!stripe) {
      return { success: false, error: 'Stripe.js未能成功加载。' };
    }

    // 5. 重定向到 Stripe Checkout 页面
    if (options.openInNewTab) {
      // 在新标签页中打开 Stripe Checkout
      const stripeWindow = window.open('', '_blank');
      
      if (stripeWindow) {
        // 在新窗口中执行重定向
        stripeWindow.document.write(`
          <html>
            <head>
              <title>go to payment...</title>
              <style>
                body {
                  display: flex;
                  justify-content: center;
                  align-items: center;
                  height: 100vh;
                  margin: 0;
                  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
                }
                .container {
                  text-align: center;
                }
              </style>
            </head>
            <body>
              <div class="container">
                <p>go to payment, please wait...</p>
              </div>
            </body>
          </html>
        `);
        
        // 将 Stripe 对象和 sessionId 传递给新窗口
        stripeWindow.sessionId = sessionId;
        
        // 在新窗口中加载 Stripe.js 并执行重定向
        const script = stripeWindow.document.createElement('script');
        script.src = 'https://js.stripe.com/v3/';
        script.onload = function() {
          const newStripe = stripeWindow.Stripe(publishableKey);
          newStripe.redirectToCheckout({ sessionId: stripeWindow.sessionId });
        };
        stripeWindow.document.head.appendChild(script);
        return { success: true, sessionId };
      }
    }
    
    // 如果无法打开新窗口或不需要新窗口，则在当前页面重定向
    const { error } = await stripe.redirectToCheckout({
      sessionId: sessionId,
    });
    
    if (error) {
      return { success: false, error: error.message || '跳转到支付页面时发生未知错误。' };
    }
    
    return { success: true, sessionId };
  } catch (err: any) {
    return { 
      success: false, 
      error: err.message || '处理支付请求时发生了一个错误。' 
    };
  }
};
