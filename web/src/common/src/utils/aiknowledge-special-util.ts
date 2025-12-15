/**
 * 轻量级节流函数
 * @param fn 需要节流的函数
 * @param wait 等待时间（毫秒）
 * @param immediate 是否立即执行
 * @param debounce 是否为防抖模式
 * @returns 节流后的函数
 */
export function liteThrottle<T extends (...args: any[]) => any>(
  fn: T,
  wait = 300,
  immediate = false,
  debounce = false
): (...args: Parameters<T>) => ReturnType<T> | undefined {
  let timer: ReturnType<typeof setTimeout> | null = null
  let lastArgs: Parameters<T> | null = null
  let result: ReturnType<T> | undefined

  const later = () => {
    timer = null
    if (lastArgs && !immediate) {
      result = fn(...lastArgs)
      lastArgs = null
    }
  }

  return function (this: any, ...args: Parameters<T>): ReturnType<T> | undefined {
    lastArgs = args

    if (timer === null) {
      if (immediate) {
        result = fn.apply(this, args)
      }
      timer = setTimeout(later, wait)
    } else if (debounce) {
      clearTimeout(timer)
      timer = setTimeout(later, wait)
    }

    return result
  }
}

/**
 * 包装函数，添加废弃警告
 * @param fn 需要标记为废弃的函数
 * @param message 废弃警告消息
 * @returns 包装后的函数
 */
export function withDeprecate<T extends (...args: any[]) => any>(
  fn: T,
  // 保留参数但不使用，以保持 API 兼容性
  _message = 'This function is deprecated and will be removed in a future version.'
): (...args: Parameters<T>) => ReturnType<T> {
  // 移除未使用的变量
  return function (this: any, ...args: Parameters<T>): ReturnType<T> {
    // 禁用废弃警告输出以减少控制台噪音
    // 原代码：
    // if (!warned) {
    //   console.warn(message)
    //   warned = true
    // }
    return fn.apply(this, args)
  }
}

/**
 * 延迟函数，返回一个 Promise，在指定时间后 resolve
 * @param ms 延迟时间（毫秒）
 * @returns Promise
 */
export function delay(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms))
}
