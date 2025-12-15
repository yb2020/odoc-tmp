import { useData } from 'vitepress';

/**
 * 获取 VitePress 的 base 路径
 * @returns VitePress 的 base 路径
 */
export function getBasePath() {
  const { site } = useData();
  return site.value.base || '/';
}

/**
 * 将路径与 VitePress 的 base 路径组合
 * @param path 需要组合的路径
 * @returns 组合后的完整路径
 */
export function withBase(path: string): string {
  const base = getBasePath();
  
  // 如果路径是绝对URL或已经包含base，直接返回
  if (path.startsWith('http') || path.startsWith(base)) {
    return path;
  }
  
  // 移除路径开头的斜杠，然后与base组合
  return base + path.replace(/^\//, '');
}

/**
 * 处理图片URL，添加base路径
 * @param imgUrl 图片URL
 * @returns 处理后的图片URL
 */
export function getImageUrl(imgUrl: string): string {
  return withBase(imgUrl);
}