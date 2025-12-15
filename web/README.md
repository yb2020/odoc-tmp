# PDF Annotate

这个项目是一个 PDF 注释工具，允许用户在 PDF 文档上添加各种注释。

## 项目结构

- `src/`: 源代码目录
  - `pdf-annotate-core/`: PDF 注释核心功能
  - `pdf-annotate-viewer/`: PDF 查看器组件
  - 其他前端代码

## 技术栈

- Node.js 20.18.1
- Vue 3.5
- Vite 6
- TypeScript
- Bun (包管理和运行时)
- Tauri 2 (桌面端应用)
- Rust (Tauri 后端)

## 安装

### 安装 Bun

#### macOS / Linux

```bash
# 安装 Bun
curl -fsSL https://bun.sh/install | bash

# 将 Bun 添加到 PATH
export PATH="$HOME/.bun/bin:$PATH"
```

#### Windows

在 Windows 上，可以通过以下方式安装 Bun：

1. **使用 PowerShell 安装**：

```powershell
# 使用 PowerShell 安装 Bun
powershell -c "irm bun.sh/install.ps1 | iex"
```

2. **使用 Windows Subsystem for Linux (WSL)**：

```bash
# 在 WSL 中安装 Bun
curl -fsSL https://bun.sh/install | bash
```

3. **使用 scoop**：  

```powershell
# 使用 scoop 安装 Bun
scoop install bun
```

安装完成后，重启终端或命令提示符以使 bun 命令生效。

### 安装 Rust (Tauri 开发必需)

Tauri 桌面端开发需要 Rust 环境。

#### macOS / Linux

```bash
# 安装 Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# 将 Cargo 添加到 PATH
echo 'export PATH="$HOME/.cargo/bin:$PATH"' >> ~/.zshrc
# 或者如果使用 bash
# echo 'export PATH="$HOME/.cargo/bin:$PATH"' >> ~/.bashrc

# 重新加载配置
source ~/.zshrc  # 或 source ~/.bashrc

# 验证安装
rustc --version
cargo --version
```

#### Windows

1. 下载并运行 [Rust 安装程序](https://www.rust-lang.org/tools/install)
2. 按照提示完成安装（需要 Visual Studio C++ Build Tools）
3. 重启终端

#### macOS 额外依赖

```bash
# 安装 Xcode Command Line Tools
xcode-select --install
```

#### Linux 额外依赖 (Ubuntu/Debian)

```bash
sudo apt update
sudo apt install libwebkit2gtk-4.1-dev libappindicator3-dev librsvg2-dev patchelf
```

### 安装项目依赖

```bash
bun install
```

## 开发

### Web 开发

```bash
bun run dev
```

### Tauri 桌面端开发

```bash
# 启动 Tauri 开发模式（会同时启动前端和 Tauri 窗口）
bun run tauri:dev
```

如果 `tauri:dev` 无法正常检测前端服务，可以使用两个终端分开运行：

```bash
# 终端 1: 启动前端
bun run dev

# 终端 2: 启动 Tauri（不启动前端）
cargo tauri dev --no-dev-server
```

## 构建

### Web 构建

```bash
bun run build
```

### Tauri 桌面端构建

```bash
# 构建当前平台的安装包
bun run tauri:build

# 构建产物位置
# macOS: src-tauri/target/release/bundle/dmg/
# Windows: src-tauri/target/release/bundle/msi/
# Linux: src-tauri/target/release/bundle/deb/
```

## 项目结构 (Tauri)

```
go-sea-web/
├── src/                          # 前端源码
│   └── utils/platform.ts         # 平台检测工具
├── src-tauri/                    # Tauri 后端
│   ├── Cargo.toml                # Rust 依赖配置
│   ├── tauri.conf.json           # Tauri 配置
│   ├── capabilities/             # 权限配置
│   ├── icons/                    # 应用图标
│   └── src/
│       ├── main.rs               # 入口
│       ├── lib.rs                # 库入口
│       └── commands.rs           # Tauri 命令
└── .github/workflows/
    └── ci-tauri.yaml             # Tauri 独立 CI
```

## 平台检测

在前端代码中可以使用 `src/utils/platform.ts` 检测当前运行环境：

```typescript
import { isTauri, isWeb, getPlatform } from '@/utils/platform'

if (isTauri()) {
  // Tauri 桌面端特有逻辑
} else {
  // Web 浏览器逻辑
}
```

## 最近更新 (2025-04-09)

1. **迁移到 Bun**: 项目已从 Yarn 迁移到 Bun 作为包管理器和运行时，提供更快的依赖安装和开发服务器启动速度。

2. **升级核心依赖**：
   - Vue 升级到 3.5.13 (从 3.3.4)
   - Vite 升级到 6.2.5 (从 2.7.2)
   - Vue Router 升级到 4.5.0
   - Pinia 升级到 3.0.1
   - 其他相关插件也已更新到最新版本

3. **配置优化**：
   - 更新了 PostCSS 配置以支持最新的插件 API
   - 优化了 Vite 配置以适应 Vite 6 的 API 变化

4. **集成 Tauri 2**：
   - 添加 Tauri 2 桌面端支持
   - 独立的 Tauri CI 工作流
   - 平台检测工具

## 常见问题

### Cargo 命令未找到

如果运行 `bun run tauri:dev` 时遇到 `cargo: command not found` 错误：

```bash
# 将 Cargo 添加到 PATH
echo 'export PATH="$HOME/.cargo/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# 验证
cargo --version
```

### Bun 命令未找到

#### macOS / Linux

如果遇到 `bun: command not found` 错误，请确保将 Bun 添加到 PATH：

```bash
echo 'export PATH="$HOME/.bun/bin:$PATH"' >> ~/.zshrc
# 或者如果使用 bash
# echo 'export PATH="$HOME/.bun/bin:$PATH"' >> ~/.bashrc
source ~/.zshrc  # 或 source ~/.bashrc
```

#### Windows

如果在 Windows 上遇到命令未找到的问题：

1. 确保 Bun 已正确安装
2. 重启命令提示符或 PowerShell
3. 检查系统环境变量是否包含 Bun 的安装路径（通常在 `%LOCALAPPDATA%\bun\bin`）

### 构建错误

如果在构建过程中遇到错误，可能是由于依赖兼容性问题。尝试清除缓存并重新安装：

```bash
rm -rf node_modules  # macOS/Linux
# 或在 Windows 上
# rmdir /s /q node_modules

bun install
```

### bun cache

如果在构建过程中遇到错误，可能是由于依赖缓存问题。尝试清除缓存并重新安装：

```bash
bun pm cache rm
bun install
```

构建使用本地 go-sea-proto
1. 修改package.json

```json
"go-sea-proto": "git+ssh://git@github.com/yb2020/go-sea-proto.git#v0.0.56" 
// 修改为
"go-sea-proto": "file:${workspaceFolder}/go-sea-proto"
例如：
"go-sea-proto": "file:D:/gitCode/go-readpaper/go-sea-proto"
"go-sea-proto": "file:/Users/jinzhi/Desktop/github/go-sea-proto"
```
2. 执行
```bash
 在 go-sea-proto 目录下执行
npm link
npm install @protobuf-ts/runtime --legacy-peer-deps

在 go-sea-web 目录下执行
npm link go-sea-proto --legacy-peer-deps
```

## git查看最后提交记录
```bash
git log --name-status -- package-lock.json
```

## 开发注意
- package.json
- vite.config.ts
修改了这两个文档，要在本地清除bun缓存，然后在本地使用bun run build看是否能打包通过，避免上ci无法编译过，又回过来找问题


## 发版本流程
目前测试阶段，一律发版本为v0.0.x-develop
1. 查看项目对应最新的tag标签
2. 将自己的开发特性分支合并到master, push到github
3. 将master分支合并到release分支，push到github
4. 在release分支上打标签，标签格式为v0.0.x-develop
5. 等待github action自动触发ci-release, 直到ci-release成功
6. 群里通知当前要发到生产环境的版本，格式如下:
```
项目：go-sea-web
版本：v0.0.x-develop
```

7. 运维修改要发到生产环境的版本，push到github
8. 等待argocd自动触发更新，或者手动点击项目的refresh按钮, 等待更新完成
9. 添加tag对应的release说明，本次更新内容

### 版本打标签命令
```bash
git tag "v0.0.2-develop"
git push origin v0.0.2-develop
```

### 版本打标签强制命令-慎用！会覆盖远程标签或者发版本失败！
```bash
git tag -f "v0.0.2-develop"
git push origin -f v0.0.2-develop
```

## 普通回滚流程
跟发版本流程一样，只不过版本继续向前推，直到发版本成功

## 紧急回滚流程
- 直接通知运维要回滚的版本，运维回滚到此版本，直接重新拉起pod即可
- 注意：受限磁盘空间的原因，并不能无限回滚到任意版本，一般的紧急回滚，仅能回滚上一版本。所以选择回滚时，一定要慎重。
- 回滚过程中，如有数据库更新，一定要注意所有的版本更新，确保数据库的修改不会影响到应用本身，以免造成数据丢失、或者生产事故！


## 许可证

MIT
