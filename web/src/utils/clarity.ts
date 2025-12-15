/**
 * Microsoft Clarity 集成工具类
 * 提供用户行为分析和会话录制功能
 */
import clarity from '@microsoft/clarity';

export class ClarityService {
  private static instance: ClarityService;
  private isInitialized = false;
  private readonly projectId = 'swn8w0ju3d';

  private constructor() {}

  public static getInstance(): ClarityService {
    if (!ClarityService.instance) {
      ClarityService.instance = new ClarityService();
    }
    return ClarityService.instance;
  }

  /**
   * 初始化 Microsoft Clarity
   */
  public init(): void {
    if (this.isInitialized) {
      console.warn('Microsoft Clarity is already initialized');
      return;
    }

    if (typeof window === 'undefined') {
      console.warn('Microsoft Clarity can only be initialized in browser environment');
      return;
    }

    try {
      // 检查是否已经初始化过
      if ((window as any).clarity) {
        console.log('Microsoft Clarity already initialized');
        return;
      }

      clarity.init(this.projectId);
      this.isInitialized = true;
      console.log(`Microsoft Clarity initialized successfully with project ID: ${this.projectId}`);
      
      // 在全局对象上标记已初始化
      (window as any).clarityInitialized = true;
    } catch (error) {
      console.error('Failed to initialize Microsoft Clarity:', error);
    }
  }

  /**
   * 设置用户标识符
   * @param userId 用户ID
   * @param sessionId 会话ID（可选）
   */
  public identify(userId: string, sessionId?: string): void {
    if (!this.isInitialized) {
      console.warn('Microsoft Clarity is not initialized. Call init() first.');
      return;
    }

    try {
      clarity.identify(userId, sessionId);
      console.log(`User identified: ${userId}${sessionId ? `, session: ${sessionId}` : ''}`);
    } catch (error) {
      console.error('Failed to identify user in Clarity:', error);
    }
  }

  /**
   * 设置自定义标签
   * @param key 标签键
   * @param value 标签值
   */
  public setTag(key: string, value: string | string[]): void {
    if (!this.isInitialized) {
      console.warn('Microsoft Clarity is not initialized. Call init() first.');
      return;
    }

    try {
      const tagValue = Array.isArray(value) ? value : [String(value)];
      clarity.setTag(key, tagValue);
      console.log(`Clarity tag set: ${key} = ${value}`);
    } catch (error) {
      console.error('Failed to set Clarity tag:', error);
    }
  }

  /**
   * 记录自定义事件（通过设置标签实现）
   * @param eventName 事件名称
   * @param properties 事件属性（可选）
   */
  public event(eventName: string, properties?: Record<string, any>): void {
    if (!this.isInitialized) {
      console.warn('Microsoft Clarity is not initialized. Call init() first.');
      return;
    }

    try {
      // 记录事件名称
      this.setTag('last_event', eventName);
      
      if (properties) {
        // 将属性设置为标签
        Object.entries(properties).forEach(([key, value]) => {
          this.setTag(`event_${eventName}_${key}`, String(value));
        });
      }
      
      console.log(`Clarity event recorded: ${eventName}`, properties || '');
    } catch (error) {
      console.error('Failed to record Clarity event:', error);
    }
  }

  /**
   * 检查 Clarity 是否已初始化
   */
  public get initialized(): boolean {
    return this.isInitialized;
  }

  /**
   * 获取项目ID
   */
  public get id(): string {
    return this.projectId;
  }
}

// 导出单例实例
export const clarityService = ClarityService.getInstance();

// 导出便捷方法
export const initClarity = () => clarityService.init();
export const identifyUser = (userId: string, sessionId?: string) => clarityService.identify(userId, sessionId);
export const setClarityTag = (key: string, value: string | string[]) => clarityService.setTag(key, value);
export const trackEvent = (eventName: string, properties?: Record<string, any>) => clarityService.event(eventName, properties);
