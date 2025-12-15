# 构建阶段
FROM oven/bun:1.2.18 AS build

# 设置工作目录
WORKDIR /app

# 定义构建参数
ARG APP_ENV=staging
ARG NODE_ENV=staging
ARG SKIP_BUILD=true
ARG SKIP_DEPS=true
ARG GITHUB_PAT
ARG BUILD_VERSION
ARG BUILD_TIME
ARG GIT_COMMIT

# 设置环境变量
ENV VITE_API_ENV=${APP_ENV}
ENV VITE_BUILD_VERSION=${BUILD_VERSION}
ENV VITE_BUILD_TIME=${BUILD_TIME}
ENV VITE_GIT_COMMIT=${GIT_COMMIT}

# 安装git
RUN apt-get update && apt-get install -y git

# 复制package.json、bun.lock和本地依赖文件以优化缓存
# 只在需要安装依赖时复制 package.json 和本地 tgz 文件
RUN if [ "$SKIP_DEPS" = "false" ]; then \
      echo "Preparing for dependency installation..."; \
    else \
      echo "Skipping dependency preparation as SKIP_DEPS=$SKIP_DEPS"; \
    fi

# 条件复制 package.json 和本地依赖
COPY --chown=node:node package.json bun.lock* ./
COPY --chown=node:node idea-*.tgz ./

# 如果需要处理 Git 依赖
RUN if [ -n "$GITHUB_PAT" ] && [ "$SKIP_DEPS" = "false" ]; then \
      # 将 SSH URL 替换为带有 GITHUB_PAT 的 HTTPS URL，只替换协议部分 \
      sed -i 's|git+ssh://git@github.com|git+https://'$GITHUB_PAT'@github.com|g' package.json && \
      echo "Successfully replaced SSH URL with HTTPS+PAT URL in package.json"; \
    fi

# 安装依赖
RUN if [ "$SKIP_DEPS" = "false" ]; then \
      echo "Installing dependencies..." && \
      bun pm cache rm && bun install; \
    else \
      echo "Skipping dependency installation as SKIP_DEPS=$SKIP_DEPS"; \
    fi

# 复制其余项目文件
COPY . .

RUN if [ "$SKIP_BUILD" = "false" ]; then \
    echo "Building application with APP_ENV=$APP_ENV, VERSION=$BUILD_VERSION" && \
    # 设置 NODE_ENV=production 以避免开发模式下的警告 \
    export NODE_ENV=production && \
    # 定义 __VUE_PROD_DEVTOOLS__ 变量，避免 vue-i18n 错误 \
    export __VUE_PROD_DEVTOOLS__=false && \
    # 使用主项目的构建命令，它会同时构建 VitePress 文档 \
    bun run build; \
  else \
    echo "Skipping build step as SKIP_BUILD=$SKIP_BUILD"; \
  fi

# 生产阶段
FROM nginx:1.25.4-alpine

# 维护者信息
LABEL maintainer="20207745@qq.com"

# 配置 nginx 日志
RUN ln -sf /dev/stderr /var/log/nginx/error.log && \
    ln -sf /dev/stdout /var/log/nginx/access.log

# 从构建阶段复制构建产物
COPY --from=build /app/dist /etc/nginx/html

# 移除默认配置
RUN rm -f /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]