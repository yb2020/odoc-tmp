import Cookies from 'js-cookie';

/**
 * 存储类型枚举
 */
export enum StorageType {
  LOCAL = 'localStorage',
  SESSION = 'sessionStorage',
  COOKIE = 'cookie'
}

/**
 * Cookie选项接口
 */
interface CookieOptions {
  expires?: number | Date;
  path?: string;
  domain?: string;
  secure?: boolean;
  sameSite?: 'strict' | 'lax' | 'none';
}

/**
 * 统一的存储管理服务
 */
export class StorageService {
  /**
   * 获取存储项
   * @param key 键名
   * @param type 存储类型
   * @returns 存储的值
   */
  static get(key: string, type: StorageType = StorageType.LOCAL): string | null {
    try {
      switch (type) {
        case StorageType.LOCAL:
          return localStorage.getItem(key);
        case StorageType.SESSION:
          return sessionStorage.getItem(key);
        case StorageType.COOKIE:
          return Cookies.get(key) || null;
        default:
          return null;
      }
    } catch (error) {
      console.error(`Error getting ${key} from ${type} storage:`, error);
      return null;
    }
  }

  /**
   * 设置存储项
   * @param key 键名
   * @param value 要存储的值
   * @param type 存储类型
   * @param options Cookie选项（仅在type为COOKIE时有效）
   */
  static set(key: string, value: string, type: StorageType = StorageType.LOCAL, options?: CookieOptions): void {
    try {
      switch (type) {
        case StorageType.LOCAL:
          localStorage.setItem(key, value);
          break;
        case StorageType.SESSION:
          sessionStorage.setItem(key, value);
          break;
        case StorageType.COOKIE:
          Cookies.set(key, value, options);
          break;
      }
    } catch (error) {
      console.error(`Error setting ${key} in ${type} storage:`, error);
    }
  }

  /**
   * 移除存储项
   * @param key 键名
   * @param type 存储类型
   */
  static remove(key: string, type: StorageType = StorageType.LOCAL): void {
    try {
      if (type === StorageType.LOCAL) {
        localStorage.removeItem(key);
      } else if (type === StorageType.SESSION) {
        sessionStorage.removeItem(key);
      } else if (type === StorageType.COOKIE) {
        Cookies.remove(key);
      }
    } catch (error) {
      console.error(`Error removing ${key} from ${type} storage:`, error);
    }
  }

  /**
   * 监听存储变化事件
   * @param callback 回调函数
   */
  static addStorageListener(callback: (event: StorageEvent) => void): void {
    window.addEventListener('storage', callback);
  }

  /**
   * 移除存储变化事件监听
   * @param callback 回调函数
   */
  static removeStorageListener(callback: (event: StorageEvent) => void): void {
    window.removeEventListener('storage', callback);
  }
}

/**
 * 创建一个简单的存储代理，用于直接访问特定类型的存储
 * @param type 存储类型
 * @returns 存储代理对象
 */
export const createStorageProxy = (type: StorageType) => {
  return {
    get: (key: string) => StorageService.get(key, type),
    set: (key: string, value: string, options?: CookieOptions) => 
      StorageService.set(key, value, type, options),
    remove: (key: string) => StorageService.remove(key, type)
  };
};

// 导出常用存储代理
export const localStore = createStorageProxy(StorageType.LOCAL);
export const sessionStore = createStorageProxy(StorageType.SESSION);
export const cookieStore = createStorageProxy(StorageType.COOKIE);
