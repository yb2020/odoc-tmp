// 这是一个 Linaria 的 mock 模块，用于在不使用 Linaria 的情况下提供兼容性
// 返回一个空函数，这样导入 css 标签的代码不会报错
export const css = (strings, ...values) => {
  // 返回一个唯一的类名，以便在运行时使用
  return '';
};

// 导出其他可能被使用的 Linaria 函数
export const styled = () => () => '';
export const cx = (...args) => args.filter(Boolean).join(' ');
export default { css, styled, cx };
