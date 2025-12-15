import path from 'path';
import { execSync } from 'child_process';
import { defineConfig } from 'vite';
import { VitePWA as pwa } from 'vite-plugin-pwa';
import vue from '@vitejs/plugin-vue';
import legacy from '@vitejs/plugin-legacy';
import crossorigin from 'vite-plugin-html-crossorigin';
import svg from 'vite-svg-loader';
import fs from 'fs/promises';
// 导入 VitePress 插件
const vitepressPlugin = require('./build/plugins/vitepress-plugin');

const IS_DEV = process.env.NODE_ENV !== 'production';
const MODE =
  process.env.VITE_API_ENV ||
  {
    develop: 'dev',
    uat: 'uat',
    production: 'prod',
  }[process.env.CICD_BUILD_ENV_NAME?.toLowerCase() || 'production'];
const HASH =
  process.env.COMMIT_ID ||
  execSync('git rev-parse HEAD').toString() ||
  Date.now().toString();

// 检查是否跳过类型检查
const SKIP_TYPECHECK = process.env.SKIP_TYPECHECK === 'true';

// https://vitejs.dev/config/
export default defineConfig({
  mode: MODE,
  base: '/',
  define: {
    'import.meta.env.REVISION': `'${HASH.trim()}'`,
    // 禁用 Linaria 的运行时检查
    'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'development'),
    '__LINARIA__': false
  },
  plugins: [
    vue({
      template: {
        compilerOptions: {
          isCustomElement: (tag) => tag.includes('custom-shadow-'),
        },
      },
    }),
    legacy({
      targets: ['defaults', 'not IE 11'],
      modernPolyfills: ['es/array/at', 'es/string/at', 'web.structured-clone'],
      renderLegacyChunks: false,
    }),
    pwa({
      // 确保在主域
      base: '/',
      strategies: 'generateSW',
      injectRegister: 'script',
      registerType: 'autoUpdate',
      workbox: {
        skipWaiting: true,
        // 增加缓存文件大小限制，默认为 2MB，现在设置为 4MB
        maximumFileSizeToCacheInBytes: 10 * 1024 * 1024, // 10MB
        // 如需预缓存更多资源
        // globPatterns: ["**\/*.{js,css,html}"],
        navigateFallbackDenylist: [
          // Do not intercept API calls or specific html files
          /^\/api\//,
          /summary\.html/,
          /client\.html/,
          /sw\.js/,
          /.+\.html/,
        ],
        runtimeCaching: [
          {
            urlPattern: /^https:\/\/nuxt\.cdn\.readpaper\.com\/.*/i,
            handler: 'CacheFirst',
            options: {
              cacheName: 'pdf-cdn-cache',
              expiration: {
                maxEntries: 10,
                maxAgeSeconds: 60 * 60 * 24 * 365, // <== 365 days
              },
              // 0的话代表缓存跨域的opaque response，所以用默认行为即可
              // cacheableResponse: {
              //   statuses: [0, 200],
              // },
            },
          },
        ],
      },
      // 如需卸载，测试环境卸载
      selfDestroying: process.env.CICD_BUILD_ENV_NAME === 'develop',
    }),
    // 添加 VitePress 插件
    vitepressPlugin({
      root: path.resolve(__dirname, 'docs'),
      base: '/docs/'
    }),
    svg({
      defaultImport: 'url',
      svgoConfig: {
        plugins: ['prefixIds'],
      },
    }),
    crossorigin({
      extensions: ['css', 'js'],
      includes: ['nuxt.cdn.readpaper.com'],
    }),
  ],
  resolve: {
    // resolve common依赖的包时强制指向 __dirname/node_modules去寻找
    // dedupe: [...Object.keys(commonPackageData.dependencies)],
    alias: {
      '@': path.join(__dirname, 'src'),
      '~': __dirname,
      '~common': path.join(__dirname, 'src/common'),
      '@common/api': path.join(__dirname, 'src/api'),
      '@common': path.join(__dirname, 'src/common/src'),
      '@idea/pdf-annotate-core': path.join(__dirname, 'src/pdf-annotate-core'),
      '@idea/pdf-annotate-viewer': path.join(__dirname, 'src/pdf-annotate-viewer'),
      '@idea/pdf-annotate-core/dist': path.join(__dirname, 'src/pdf-annotate-core'),
      '@idea/pdf-annotate-viewer/dist': path.join(__dirname, 'src/pdf-annotate-viewer'),
      // 添加 Linaria 的 mock 模块，返回一个空函数
      '@home/utils/citation': path.join(__dirname, 'src/utils/citation.ts'),
      'linaria': path.resolve(__dirname, 'src/mocks/linaria-mock.js'),
      '@linaria/core': path.resolve(__dirname, 'src/mocks/linaria-mock.js'),
      // Use our custom Vue 3 implementation of vue-intersect
      'vue-intersect': path.resolve(__dirname, 'src/shims/vue-intersect.js'),
    },
    extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue']
  },
  css: {
    preprocessorOptions: {
      less: {
        javascriptEnabled: true,
      },
    },
  },
  server: {
    fs: {
      allow: [path.join(__dirname, '../../')],
    },
    host: true,
    port: 3000,
    proxy: {
      // 对 /api 规则进行修改
      '/api': {
        target: process.env.API_TARGET || 'http://localhost:8081',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api'),
        /**
         * 添加 configure 钩子来修改代理请求头。
         * 这是解决问题的关键。
         */
        configure: (proxy, options) => {
          // 监听 'proxyReq' 事件，这个事件在请求被发送到后端目标前触发
          proxy.on('proxyReq', (proxyReq, req, res) => {
            // `req` 是来自浏览器的原始请求
            // `proxyReq` 是即将被发送到后端的请求
            if (req.headers.host) {
              // 将浏览器原始请求的 host 设置到 X-Forwarded-Host Header 中
              proxyReq.setHeader('X-Forwarded-Host', req.headers.host);
            }
          });
        }
      },
      // 对 /report/collection_tracking0 规则也进行同样的修改，保持一致性
      '/report/collection_tracking0': {
        target: process.env.API_TARGET || 'http://localhost:8081',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/report\/collection_tracking0/, '/report/collection_tracking0'),
        /**
         * 添加 configure 钩子来修改代理请求头。
         */
        configure: (proxy, options) => {
          proxy.on('proxyReq', (proxyReq, req, res) => {
            if (req.headers.host) {
              proxyReq.setHeader('X-Forwarded-Host', req.headers.host);
            }
          });
        }
      },
    },
  },
  build: {
    emptyOutDir: true,
    minify: true,
    rollupOptions: {
      input: {
        index: path.resolve(__dirname, 'index.html'),
        client: path.resolve(__dirname, 'client.html'),
        summary: path.resolve(__dirname, 'summary.html'),
        aibeans: path.resolve(__dirname, 'aibeans.html'),
      },
      output: {
        // 改善代码分割策略
        manualChunks: {
          // 将Vue相关库分离到单独的chunk
          'vue-vendor': ['vue', 'vue-router', 'vuex', 'vue-i18n'],
          // 将Ant Design Vue分离到单独的chunk
          'antd-vendor': ['ant-design-vue'],
          // 将大型工具库分离
          'utils-vendor': ['axios', 'lodash', 'moment'],
          // 将PDF相关功能分离
          'pdf-vendor': [],
          // 将AI相关功能分离
          'ai-vendor': []
        },
        // 优化chunk文件名
        chunkFileNames: 'js/[name]-[hash].js',
        entryFileNames: 'js/[name]-[hash].js',
        assetFileNames: (assetInfo) => {
          if (!assetInfo.name) return 'assets/[name]-[hash][extname]';
          
          const info = assetInfo.name.split('.');
          const ext = info[info.length - 1];
          if (/\.(mp4|webm|ogg|mp3|wav|flac|aac)(\?.*)?$/i.test(assetInfo.name)) {
            return `media/[name]-[hash].${ext}`;
          } else if (/\.(png|jpe?g|gif|svg)(\?.*)?$/i.test(assetInfo.name)) {
            return `images/[name]-[hash].${ext}`;
          } else if (ext === 'css') {
            return `css/[name]-[hash].${ext}`;
          }
          return `assets/[name]-[hash].${ext}`;
        }
      }
    },
    // 设置chunk大小警告阈值
    chunkSizeWarningLimit: 1000,
    // 启用源码映射用于调试（生产环境可关闭）
    sourcemap: IS_DEV
  },
  optimizeDeps: {
    exclude: ['linaria', '@linaria/core'],
    esbuildOptions: {
      plugins: [
        {
          name: 'fix-proto-imports',
          setup(build) {
            // 匹配所有 go-sea-proto 生成的 TypeScript 文件
            build.onLoad({ filter: /go-sea-proto\/gen\/ts\/.*\.ts$/ }, async (args) => {
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
  }
});
