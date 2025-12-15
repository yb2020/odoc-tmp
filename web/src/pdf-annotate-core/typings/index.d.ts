declare module 'create-stylesheet' {
  function createStyleSheet(style: Record<string, Partial<CSSStyleDeclaration>>): HTMLStyleElement;
  export default createStyleSheet;
}

declare module 'element-visible';