// docs/.vitepress/env.d.ts
/// <reference types="vitepress/client" />

declare module '*.vue' {
    import type { DefineComponent } from 'vue'
    const component: DefineComponent<{}, {}, any>
    export default component
  }