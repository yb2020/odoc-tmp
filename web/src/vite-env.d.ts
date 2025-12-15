/// <reference types="vite/client" />

// 删除自定义的 ImportMeta 接口，避免与 Vite 的类型定义冲突

interface I18nMessageType {
  [key: string]: any;
}
