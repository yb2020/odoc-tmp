/**
 * 文件选择工具
 * 根据运行环境自动选择使用 Tauri API 或浏览器原生方式
 * 
 * 注意：Tauri 模式需要安装 @tauri-apps/plugin-dialog
 * 运行: bun add @tauri-apps/plugin-dialog
 */

import { isInTauri } from '@/util/env';

/**
 * 文件选择结果
 */
export interface FileSelectResult {
  /** 文件对象（浏览器模式下可用） */
  file?: File;
  /** 文件本地路径（Tauri 模式下可用） */
  localPath?: string;
  /** 文件名 */
  fileName: string;
  /** 是否为本地路径模式 */
  isLocalPathMode: boolean;
}

/**
 * 文件选择选项
 */
export interface FileSelectOptions {
  /** 是否允许多选 */
  multiple?: boolean;
  /** 是否选择文件夹 */
  directory?: boolean;
  /** 文件过滤器 */
  filters?: {
    name: string;
    extensions: string[];
  }[];
}

/**
 * 默认 PDF 过滤器
 */
const DEFAULT_PDF_FILTER = [{ name: 'PDF', extensions: ['pdf'] }];

/**
 * 动态获取 Tauri dialog API
 * 使用动态导入避免在非 Tauri 环境下报错
 */
async function getTauriDialogApi() {
  try {
    return await import('@tauri-apps/plugin-dialog');
  } catch (error) {
    console.warn('Tauri dialog plugin not available:', error);
    return null;
  }
}

/**
 * 使用 Tauri 文件选择对话框
 */
async function selectFilesWithTauri(options: FileSelectOptions): Promise<FileSelectResult[]> {
  const dialogApi = await getTauriDialogApi();
  
  if (!dialogApi) {
    console.warn('Tauri dialog API not available, falling back to browser mode');
    return selectFilesWithBrowser(options);
  }
  
  try {
    const selected = await dialogApi.open({
      multiple: options.multiple ?? true,
      directory: options.directory ?? false,
      filters: options.filters ?? DEFAULT_PDF_FILTER,
    });
    
    if (!selected) {
      return [];
    }
    
    // 统一处理为数组
    const paths = Array.isArray(selected) ? selected : [selected];
    
    return paths.map(path => ({
      localPath: path,
      fileName: extractFileName(path),
      isLocalPathMode: true,
    }));
  } catch (error) {
    console.error('Tauri file selection failed:', error);
    throw error;
  }
}

/**
 * 使用浏览器原生文件选择
 */
async function selectFilesWithBrowser(options: FileSelectOptions): Promise<FileSelectResult[]> {
  return new Promise((resolve, reject) => {
    const input = document.createElement('input');
    input.type = 'file';
    input.multiple = options.multiple ?? true;
    input.accept = '.pdf';
    
    if (options.directory) {
      input.webkitdirectory = true;
    }
    
    input.onchange = (event) => {
      const files = (event.target as HTMLInputElement).files;
      if (!files || files.length === 0) {
        resolve([]);
        return;
      }
      
      const results: FileSelectResult[] = Array.from(files)
        .filter(file => file.type === 'application/pdf' || file.name.toLowerCase().endsWith('.pdf'))
        .map(file => ({
          file,
          fileName: file.name,
          isLocalPathMode: false,
        }));
      
      resolve(results);
    };
    
    input.onerror = () => {
      reject(new Error('File selection failed'));
    };
    
    // 处理用户取消选择的情况
    input.oncancel = () => {
      resolve([]);
    };
    
    input.click();
  });
}

/**
 * 从路径中提取文件名
 */
function extractFileName(path: string): string {
  // 处理 Windows 和 Unix 风格的路径
  const parts = path.split(/[/\\]/);
  return parts[parts.length - 1] || path;
}

/**
 * 选择文件（自动根据环境选择方式）
 * @param options 选择选项
 * @returns 文件选择结果数组
 */
export async function selectFiles(options: FileSelectOptions = {}): Promise<FileSelectResult[]> {
  if (isInTauri()) {
    return selectFilesWithTauri(options);
  }
  return selectFilesWithBrowser(options);
}

/**
 * 选择单个 PDF 文件
 * @returns 文件选择结果，如果取消则返回 null
 */
export async function selectPdfFile(): Promise<FileSelectResult | null> {
  const results = await selectFiles({
    multiple: false,
    filters: DEFAULT_PDF_FILTER,
  });
  return results.length > 0 ? results[0] : null;
}

/**
 * 选择多个 PDF 文件
 * @returns 文件选择结果数组
 */
export async function selectPdfFiles(): Promise<FileSelectResult[]> {
  return selectFiles({
    multiple: true,
    filters: DEFAULT_PDF_FILTER,
  });
}

/**
 * 选择包含 PDF 的文件夹
 * @returns 文件选择结果数组
 */
export async function selectPdfFolder(): Promise<FileSelectResult[]> {
  return selectFiles({
    multiple: true,
    directory: true,
    filters: DEFAULT_PDF_FILTER,
  });
}

/**
 * 检查当前是否为本地路径模式（Tauri 环境）
 */
export function isLocalPathMode(): boolean {
  return isInTauri();
}
