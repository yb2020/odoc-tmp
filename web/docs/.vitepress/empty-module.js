// 空模块，用于替代 @idea/aiknowledge-report 和其他浏览器特有功能
// 避免在 SSR 过程中访问 window 和 document 对象

// 实现 Clock 类
export class Clock {
  constructor() {}
  // 实现必要的方法
  start() {}
  stop() {}
  reset() {}
  getTime() { return 0; }
}

// 实现 reporter 对象
const reporter = {
  config: {
    api: '/report/collection_tracking0'
  },
  report: () => Promise.resolve(),
};

// 从错误堆栈中发现的函数
// Cr=()=>!!window.navigator.userAgent.match(/readpaper/i)
export const Cr = () => false; // 安全版本，始终返回 false

// MO=()=>Cr()&&document.referrer===""
export const MO = () => false; // 安全版本

// Mb=()=>window.location.origin
export const Mb = () => ''; // 安全版本，返回空字符串

// ix=()=>window.location.hostname
export const ix = () => ''; // 安全版本，返回空字符串

// 添加可能被使用的其他浏览器相关函数
export const eC = () => {}; // 在错误代码中看到的另一个函数

// 默认导出 reporter
export default reporter;

// 其他可能被导入的函数和类
export const someFunction = () => {};
export const someOtherFunction = () => {};

// 导出一个包含所有函数的对象，以适应不同的导入方式
export const ti = {
  config: {
    api: '/report/collection_tracking0'
  },
  report: () => Promise.resolve()
};
