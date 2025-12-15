import { defineConfig } from 'vitepress'
import { fileURLToPath } from 'url'
import { resolve } from 'path'
import fs from 'fs/promises'
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { languageEnumToStandard } from '../../src/shared/language/service';


// 解析路径
const __dirname = fileURLToPath(new URL('.', import.meta.url))
const projectRoot = resolve(__dirname, '../..')

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "ODOC.AI",
  description: "a open doc for ai ",
  base: '/docs', 
  cleanUrls: true,
  lastUpdated: true,
  // 禁用 MPA 模式，确保客户端组件能够正确渲染
  mpa: false,
  head: [
    [
      'link',
      {
        rel: 'icon',
        // href: '/favicon.ico'
      }
    ]
  ],
  locales: {
    root: {
      label: 'English',
      lang: languageEnumToStandard(Language.EN_US), // 使用标准格式 'en-US'
      // title: 'VitePress Site',
      // description: 'A VitePress Site',
      themeConfig: {
        // logo: '/favicon.ico',
        outline: {
          level: [2, 3], // This is the key setting. It tells VitePress to show both <h2> and <h3> headings.
          label: 'On this page' // This is the title text that appears above the outline links.
        },
        nav: [
          { text: 'Home', link: '/' },
          { text: 'My Library', link: '/workbench_placeholder', target: '_self' },
          { text: 'Guide', link: '/guide' },
          { text: 'Pricing', link: '/pricing' },
        ],
        sidebar: [
          {
            text: 'Introduction',
            collapsible: true, // 开启折叠
            collapsed: false,   // 默认展开
            items: [
              { text: 'What is ODoc AI', link: '/guide' }, // Assumes English content is at /guide
              { text: 'Getting Started', link: '/guide/introduction/start' },
            ]
          },
          {
            text: 'Feature Guide',
            collapsible: true, 
            collapsed: false,   
            items: [
              { text: 'Preface', link: '/guide/chapter/preface' },
              { text: '1. Use Manual', link: '/guide/chapter/useManual' },
              { text: '2. Home Top Toolbar Introduction', link: '/guide/chapter/homeTopToolbarMt' },
              { text: '3. Document Management', link: '/guide/chapter/documentMt' },
              { text: '4. Note Management', link: '/guide/chapter/noteMt' },
              { text: '5. Tag Management', link: '/guide/chapter/tagMt' },
              { text: '6. Translation Management', link: '/guide/chapter/translationMt' },
              { text: '7. Note Function', link: '/guide/chapter/noteFc' },
              { text: '8. Read page toolbar introduction', link: '/guide/chapter/toolbarMt' },
              { text: '9. AI assisted reading', link: '/guide/chapter/aiMt' }
            ]
          },
          {
            text: 'Team',
            collapsible: true, 
            collapsed: false,   
            items: [
              { text: 'About Us', link: '/guide/team/about' },
              { text: 'Updates', link: '/guide/team/update' },
              { text: 'Join Us', link: '/guide/team/join' }
            ]
          },
          {
            text: 'product',
            collapsible: true, 
            collapsed: false,   
            items: [
              { text: 'features', link: '/guide/product/features' },
              { text: 'usage', link: '/guide/product/usage' },
              { text: 'userAgreement', link: '/guide/product/userAgreement' },
              { text: 'privacy', link: '/guide/product/privacy' }
            ]
          }
        ],
        // editLinkText: 'Edit this page on GitHub',
        // lastUpdatedText: 'Last Updated',
      }
    },
    zh: {
      label: '简体中文',
      lang: languageEnumToStandard(Language.ZH_CN), // 使用标准格式 'zh-CN'
      link: '/zh/',
      // title: 'VitePress 站点',
      // description: '一个 VitePress 站点',
      themeConfig: {
        // logo: '/favicon.ico',
        outline: {
          level: [2, 3], // This is the key setting. It tells VitePress to show both <h2> and <h3> headings.
          label: '页面导航' // This is the title text that appears above the outline links.
        },
        nav: [
          { text: '首页', link: '/zh/' },
          { text: '我的文库', link: '/workbench_placeholder', target: '_self' },
          { text: '指南', link: '/zh/guide' },
          { text: '价格', link: '/zh/pricing' },
        ],
        sidebar: [
          {
            text: '简介',
            collapsible: true,
            collapsed: false, 
            items: [
              { text: '什么是ODoc Ai', link: '/zh/guide' },
              { text: '快速开始', link: '/zh/guide/introduction/start' },
            ]
          },
          {
            text: '功能介绍',
            collapsible: true, 
            collapsed: false,   
            items: [
              { text: '前言', link: '/zh/guide/chapter/preface' },
              { text: '1. 使用手册', link: '/zh/guide/chapter/useManual' },
              { text: '2. 首页顶部工具栏介绍', link: '/zh/guide/chapter/homeTopToolbarMt' },
              { text: '3. 文献管理', link: '/zh/guide/chapter/documentMt' },
              { text: '4. 笔记管理', link: '/zh/guide/chapter/noteMt' },
              { text: '5. 标签管理', link: '/zh/guide/chapter/tagMt' },
              { text: '6. 翻译管理', link: '/zh/guide/chapter/translationMt' },
              { text: '7. 笔记功能', link: '/zh/guide/chapter/noteFc' },
              { text: '8. 阅读页面工具栏介绍', link: '/zh/guide/chapter/toolbarMt' },
              { text: '9. AI辅助阅读', link: '/zh/guide/chapter/aiMt' }
            ]
          },
          {
            text: '团队',
            collapsible: true, 
            collapsed: false,   
            items: [
              { text: '关于我们', link: '/zh/guide/team/about' },
              { text: '更新', link: '/zh/guide/team/update' },
              { text: '加入我们', link: '/zh/guide/team/join' }
            ]
          },
          {
            text: '产品',
            collapsible: true, 
            collapsed: false,   
            items: [
              { text: '功能介绍', link: '/zh/guide/product/features' },
              { text: '使用说明', link: '/zh/guide/product/usage' },
              { text: '用户协议', link: '/zh/guide/product/userAgreement' },
              { text: '隐私政策', link: '/zh/guide/product/privacy' }
            ]
          }
        ],
        // editLinkText: '在 GitHub 上编辑此页',
        // lastUpdatedText: '最后更新',
      }
    }
  },
  // Global themeConfig for non-localized settings
  themeConfig: {
    // socialLinks: [
    //   { icon: 'github', link: 'https://github.com/vuejs/vitepress' }
    // ]
    // Other global theme settings like search, algolia etc. can go here
  },
  
  // Vite 配置，包括别名
  vite: {
    resolve: {
      alias: {
        // '~' resolves to the project root (e.g., /Users/yibing/go-sea-web)
        '~': resolve(projectRoot),
        // '@' resolves to the main project's src directory (e.g., /Users/yibing/go-sea-web/src)
        '@': resolve(projectRoot, 'src'),
        // '@common' resolves to the main project's src/common/src directory
        '@common': resolve(projectRoot, 'src/common/src'),
        // '~common' resolves to the main project's src/common directory
        '~common': resolve(projectRoot, 'src/common'),
        // 创建一个空的模块别名，避免 SSR 过程中访问 window
        '@idea/aiknowledge-report': resolve(__dirname, './empty-module.js'),
        '@common/locals/i18n': resolve(__dirname, 'src/common/locals/i18n'),
        '@common/locals/cookie': resolve(__dirname, 'src/common/locals/cookie'),
        '@common/hooks/useI18nLocal': resolve(__dirname, 'src/common/hooks/useI18nLocal'),
        '@common/theme': resolve(__dirname, 'src/common/theme'),
      }
    },
    // 确保能够正确处理主项目中的模块
    optimizeDeps: {
      include: ['pinia', 'vue', 'vue-i18n'],
      esbuildOptions: {
        plugins: [
          {
            name: 'fix-proto-imports',
            setup(build: any) {
              // 匹配所有 go-sea-proto 生成的 TypeScript 文件
              build.onLoad({ filter: /go-sea-proto\/gen\/ts\/.*\.ts$/ }, async (args: any) => {
                const content = await fs.readFile(args.path, 'utf8');
                
                // 检查文件是否已经有默认导出
                if (content.includes('export default')) {
                  return { contents: content, loader: 'ts' };
                }
                
                // 提取所有导出的枚举和接口名称
                const exportedNames = [];
                const exportRegex = /export (enum|interface|type|class) (\w+)/g;
                let match;
                
                while ((match = exportRegex.exec(content)) !== null) {
                  exportedNames.push(match[2]);
                }
                
                if (exportedNames.length === 0) {
                  return { contents: content, loader: 'ts' };
                }
                
                // 添加默认导出
                const defaultExport = `\n// 添加默认导出\nexport default { ${exportedNames.join(', ')} };\n`;
                
                return {
                  contents: content + defaultExport,
                  loader: 'ts',
                };
              });
            },
          },
        ],
      }
    },
    // SSR 相关配置
    ssr: {
      // 避免将这些模块外部化，这样它们会被打包进 SSR 包中
      noExternal: [
        '@idea/aiknowledge-report',
        'ant-design-vue',
        'pinia',
        '@common/i18n',
        '@common/theme',
        'vue',
        'vue-i18n',
        '@intlify/message-compiler',
      ]
    },
    define: {
      __VUE_PROD_DEVTOOLS__: false,
    },
  }
})
