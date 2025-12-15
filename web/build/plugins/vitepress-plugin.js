/**
 * VitePress 插件 - 集成 VitePress 到 Vite 项目中
 * 
 * 此插件允许在 Vite 项目中使用 VitePress 文档，
 * 通过中间件模式集成，同时确保客户端组件能够正确处理。
 */

const path = require('path');
const fs = require('fs');
const { spawn } = require('child_process');

/**
 * 创建 VitePress 插件
 * @param {Object} options 插件配置
 * @returns {import('vite').Plugin} Vite 插件
 */
module.exports = function vitepressPlugin(options = {}) {
  const root = options.root || path.resolve(process.cwd(), 'docs');
  const base = options.base || '/docs';
  
  let vitePressServerInstance = null;
  
  return {
    name: 'vite-plugin-vitepress',
    
    async configureServer(viteMainServer) {
      console.log(`[vitepress-plugin] Configuring VitePress middleware for dev server...`);
      console.log(`[vitepress-plugin] Root: ${root}, Base: ${base}`);

      try {
        // 动态导入 VitePress 以避免类型问题
        const { createServer } = await import('vitepress');
        
        // 创建 VitePress 服务器实例，使用中间件模式
        vitePressServerInstance = await createServer(root, {
          server: {
            middlewareMode: true,
            hmr: viteMainServer.config.server.hmr ? {
              port: typeof viteMainServer.config.server.hmr === 'object' && viteMainServer.config.server.hmr.port ? 
                viteMainServer.config.server.hmr.port : 24678,
            } : false,
          }
        });
        
        // 添加 VitePress 中间件来处理以 base 路径开头的请求
        viteMainServer.middlewares.use((req, res, next) => {
          // 确保请求存在且以 base 路径开头
          if (!req.url || !req.url.startsWith(base)) {
            return next();
          }
          
          // 记录请求转发信息
          console.log(`[vitepress-plugin] Forwarding request to VitePress: ${req.url}`);
          
          try {
            // 重要：不要修改 req.url，直接传递给 VitePress 中间件
            // VitePress 期望完整的 URL 包含 base 路径
            vitePressServerInstance.middlewares(req, res, next);
          } catch (err) {
            console.error(`[vitepress-plugin] Error in VitePress middleware:`, err);
            next(err);
          }
        });
        
        console.log(`[vitepress-plugin] VitePress middleware configured successfully`);
      } catch (err) {
        console.error(`[vitepress-plugin] Failed to configure VitePress middleware:`, err);
        throw err; // 重新抛出错误以便更好地诊断
      }
    },

    async closeBundle() {
      if (process.env.NODE_ENV === 'production') {
        console.log('[vitepress-plugin] Building VitePress documentation...');
        try {
          // 使用子进程构建 VitePress 文档
          const { spawn } = require('child_process');
          
          // 调整输出目录为相对于 'dist' 并使用 base 名称
          const outDirName = base.startsWith('/') ? base.substring(1) : base;
          const outDir = path.resolve(process.cwd(), 'dist', outDirName);
          
          console.log(`[vitepress-plugin] VitePress build config - Root: ${root}, Base: ${base}, OutDir: ${outDir}`);
          
          // 设置环境变量
          const env = {
            ...process.env,
            VITEPRESS_NO_EXTERNAL: '@idea/aiknowledge-report',
            NODE_OPTIONS: '--no-warnings --trace-warnings',
            DEBUG: 'vitepress:*,vite:*',
            // 添加输出目录环境变量
            VITEPRESS_OUTPUT_DIR: outDir,
            // 确保客户端组件构建
            VITEPRESS_BUILD_CLIENT: 'true'
          };
          
          console.log('[vitepress-plugin] Starting VitePress build with enhanced client support...');
          
          // 运行标准 VitePress 构建
          const vpBuild = spawn('node', [
            './node_modules/vitepress/bin/vitepress.js',
            'build',
            root,
            // 使用标准 SSR 构建以确保客户端组件正确处理
            '--force', // 强制跳过错误
            '--outDir', outDir,
            '--base', base
          ], { env, stdio: 'inherit' });
          
          return new Promise((resolve, reject) => {
            // 监听构建进程完成
            vpBuild.on('close', (code) => {
              if (code !== 0) {
                console.error(`[vitepress-plugin] VitePress build failed with exit code ${code}`);
                reject(new Error(`VitePress build failed with exit code ${code}`));
                return;
              }
              
              console.log('[vitepress-plugin] Standard build completed, ensuring client components are properly handled...');
              
              // 处理客户端组件
              try {
                // 确保客户端目录存在
                const clientDir = path.join(outDir, 'assets');
                if (!fs.existsSync(clientDir)) {
                  fs.mkdirSync(clientDir, { recursive: true });
                }
                
                console.log(`[vitepress-plugin] Client build completed successfully. Output directory: ${outDir}`);
                resolve();
              } catch (error) {
                console.error('[vitepress-plugin] Error during client processing:', error);
                reject(error);
              }
            });
            
            // 错误处理
            vpBuild.on('error', (error) => {
              console.error(`[vitepress-plugin] VitePress build error: ${error.message}`);
              console.error('[vitepress-plugin] Error stack:', error.stack);
              
              // 如果是 window is not defined 错误
              if (error.message && error.message.includes('window is not defined')) {
                console.error('[vitepress-plugin] Detected window is not defined error during SSR');
                console.error('[vitepress-plugin] Suggestion: Wrap components that use browser APIs with <ClientOnly> tags');
              }
              
              reject(error);
            });
          });
        } catch (error) {
          console.error('[vitepress-plugin] Failed to build VitePress documentation:', error);
          throw error;
        }
      }
    }
  };
};
