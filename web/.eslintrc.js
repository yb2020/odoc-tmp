module.exports = {
  root: true,
  env: {
    browser: true,
    node: true,
    es2021: true,
  },
  extends: [
    'plugin:vue/vue3-recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:promise/recommended',
  ],
  parser: 'vue-eslint-parser',
  parserOptions: {
    ecmaVersion: 2021,
    parser: '@typescript-eslint/parser',
    sourceType: 'module',
    ecmaFeatures: {
      jsx: true,
    },
  },
  plugins: ['vue', '@typescript-eslint', 'promise'],
  rules: {
    // Vue 规则
    'vue/multi-word-component-names': 'off', // 允许单词组件名
    'vue/no-multiple-template-root': 'off', // 允许多根节点模板
    'vue/attribute-hyphenation': 'off', // 允许驼峰属性名
    'vue/v-on-event-hyphenation': 'off', // 允许驼峰事件名
    'vue/no-v-html': 'off', // 允许使用 v-html
    'vue/require-prop-types': 'warn', // 将属性类型检查降级为警告
    'vue/valid-v-for': 'warn', // 将 v-for 检查降级为警告
    'vue/return-in-computed-property': 'warn', // 将计算属性返回值检查降级为警告
    'vue/require-valid-default-prop': 'warn', // 将 prop 默认值检查降级为警告
    'vue/require-v-for-key': 'warn', // 将 v-for key 检查降级为警告
    'vue/no-dupe-keys': 'warn', // 将重复键检查降级为警告
    
    // TypeScript 规则
    '@typescript-eslint/no-explicit-any': 'warn', // 将 any 类型检查降级为警告
    '@typescript-eslint/no-unused-vars': 'warn', // 将未使用变量检查降级为警告
    '@typescript-eslint/ban-types': 'warn', // 将类型检查降级为警告
    '@typescript-eslint/ban-ts-comment': 'warn', // 将 ts-ignore 注释检查降级为警告
    
    // Promise 规则
    'promise/always-return': 'warn', // 将 Promise 返回值检查降级为警告
    
    // 其他规则
    'comma-dangle': [
      'warn',
      {
        arrays: 'only-multiline',
        objects: 'always-multiline',
        imports: 'only-multiline',
        exports: 'only-multiline',
        functions: 'never',
      }
    ],
  },
  ignorePatterns: ['node_modules/**', 'dist/**', '.git/**', '**/*.min.js'],
};
