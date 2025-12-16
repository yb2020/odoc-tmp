# Go Sea Proto

Protocol Buffers 定义库，包含各种服务的协议定义和生成的代码。

## 目录结构

- `definitions/`: 协议定义文件 (.proto)
- `gen/`: 生成的代码
  - `go/`: Go 语言生成代码
  - `ts/`: TypeScript 生成代码
  - `py/`: Python 生成代码

## 使用方法

### 本地开发前必看！！！

1. 使用 `make local` 命令编译 proto 文件到 `local-gen` 目录，查看是否正常编译。
2. 一定要使用 `make local` 而不是 `make all`，因为 `make all` 会编译到 `gen` 目录，而 `gen` 目录被 Git 忽略。
3. 提交前要使用 `make local` 检查proto是否正常编译，避免流水线无法正常运行，导致需要人为介入。

### windows 命令
1. 通过 https://www.msys2.org/ 下载MSYS2工具，使用默认地址C:\msys64安装。
2. 安装成功之后在出现的窗口中执行pacman -S mingw-w64-x86_64-make命令安装make工具
3. 将C:\msys64\mingw64\bin添加环境变量path中
4. 通过mingw32-make --version验证是否安装成功
5. 通过 https://github.com/protocolbuffers/protobuf/releases 安装 Protocol Buffers 编译器 (protoc) 下载 Windows 版本，比如 protoc-25.1-win64.zip
6. 将安装的目录&path/bin目录添加到path环境变量中，例如C:\protoc\bin
7. 安装Go 的 protoc 插件：
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    
在windows中，没有make命令，只有mingw32-make命令，以上安装成功之后就可以使用mingw32-make命令了，命令使用如下：

### windows批处理文件

项目包含以下批处理文件，用于处理 protobuf 生成的文件：

- `move_go_files.bat`: 将 `gen/go/definitions` 目录下生成的 Go 文件移动到正确的模块目录结构中
- `move_ts_files.bat`: 将 `gen/ts-sea-proto` 目录下生成的 TypeScript 文件移动到正确的模块目录结构中
- `move_local_go_files.bat`: 将 `local-gen/go/definitions` 目录下生成的 Go 文件移动到正确的模块目录结构中
- `move_local_ts_files.bat`: 将 `local-gen/ts-sea-proto` 目录下生成的 TypeScript 文件移动到正确的模块目录结构中
- `check_ts_plugin.bat`: 检查 TypeScript protoc 插件是否已安装

这些批处理文件由 MakefileForWindows 自动调用，通常不需要手动执行。

```bash
# 本地清空
mingw32-make -f makefileForWindows clean
# 本地测试
mingw32-make -f makefileForWindows local
```

### git设置大小写敏感

```bash
# 项目
git config core.ignorecase false

# 全局
git config --global core.ignorecase false
```

### Go 项目

1. 配置环境以访问私有仓库：

#### Linux/Mac 环境

```bash
# 告诉 Go 直接从源获取，不使用代理
export GOPRIVATE=github.com/yb2020/odoc/proto
# 永久设置到zsh
export GOPRIVATE=github.com/yb2020/odoc/proto
source ~/.zshrc

# 配置 Git 使用 SSH 而不是 HTTPS（只针对这个特定仓库）
git config --global url."git@github.com:yb2020/go-sea-proto.git".insteadOf "https://github.com/yb2020/odoc/proto"
```

#### Windows 环境

```powershell
# PowerShell 中设置环境变量
$env:GOPRIVATE = "github.com/yb2020/odoc/proto"
# 永久设置（需要管理员权限）
[Environment]::SetEnvironmentVariable("GOPRIVATE", "github.com/yb2020/odoc/proto", "User")

# Git 配置与 Linux/Mac 相同
git config --global url."git@github.com:yb2020/go-sea-proto.git".insteadOf "https://github.com/yb2020/odoc/proto"
```

2. 在 `go.mod` 文件中添加依赖：

```go
require github.com/yb2020/odoc/proto v0.0.3 // 使用最新版本
```

3. 在代码中导入：

```go
import (
    user "github.com/yb2020/odoc/proto/gen/go/user"
    common "github.com/yb2020/odoc/proto/gen/go/common"
)
```

### TypeScript 项目

1. 在 `package.json` 中添加依赖：

```json
{
  "dependencies": {
    "go-sea-proto": "git+ssh://git@github.com/yb2020/odoc/proto.git#v0.0.6",
  }
}
```

2. 在代码中导入：

```typescript
import { User } from 'go-sea-proto/gen/ts/User/User';
```

### Python 项目

1. 配置环境以访问私有仓库：

#### Linux/Mac 环境

```bash
# 配置 Git 使用 SSH 而不是 HTTPS（只针对这个特定仓库）
git config --global url."git@github.com:yb2020/go-sea-proto.git".insteadOf "https://github.com/yb2020/odoc/proto"
```

#### Windows 环境

```powershell
# Git 配置与 Linux/Mac 相同
git config --global url."git@github.com:yb2020/go-sea-proto.git".insteadOf "https://github.com/yb2020/odoc/proto"
```

2. 安装依赖（使用 pip）：

```bash
# 安装特定版本
pip install git+ssh://git@github.com/yb2020/odoc/proto.git@v0.0.137

# 或者在 requirements.txt 中添加
# go-sea-proto @ git+ssh://git@github.com/yb2020/odoc/proto.git@v0.0.137
```

3. 在代码中导入：

```python
from user import user_pb2
from common import common_pb2
```

## 更新协议

1. 修改 `definitions/` 目录下的 `.proto` 文件
2. 提交并推送到 `master` 分支
3. GitHub Actions 会自动生成代码并创建新的版本标签

## 获取更新

### Go 项目

```bash
go get -u github.com/yb2020/odoc/proto
```

### TypeScript 项目

更新 `package.json` 中的版本号，然后运行：

```bash
bun install
```

### Python 项目

```bash
pip install --upgrade git+ssh://git@github.com/yb2020/odoc/proto.git@v0.0.137
``` 

前端通过git拉包拉不下来时，要执行以下命令
```bash
bun pm cache rm
```


### 如何在windows环境下使用本地proto文件

1. 首先增加package-lock.json文件,内容如下：
{
  "name": "go-sea-proto",
  "version": "0.0.70",
  "lockfileVersion": 3,
  "requires": true,
  "packages": {
    "": {
      "name": "go-sea-proto",
      "version": "0.0.70",
      "dependencies": {
        "@protobuf-ts/runtime": "^2.11.0"
      }
    },
    "node_modules/@protobuf-ts/runtime": {
      "version": "2.11.0",
      "resolved": "https://registry.npmmirror.com/@protobuf-ts/runtime/-/runtime-2.11.0.tgz",
      "integrity": "sha512-DfpRpUiNvPC3Kj48CmlU4HaIEY1Myh++PIumMmohBAk8/k0d2CkxYxJfPyUAxfuUfl97F4AvuCu1gXmfOG7OJQ==",
      "license": "(Apache-2.0 AND BSD-3-Clause)"
    }
  }
}
2.  然后在当前项目中依次运行如下命令: 执行完命令之后会在当前目录下生成一个node_modules文件夹  

```bash
npm link 
```
```bash
npm install @protobuf-ts/runtime
```

3. 然后在你的引用proto协议的项目中执行

```bash
npm link go-sea-proto --legacy-peer-deps
```

### 工具脚本文件说明

1. delete_windows_local_config.bat
  删除windows本地配置文件,清理生成的文件，方便git提交
2. gen_windows_local_web_config.bat
  生成windows本地配置文件,主要应用于前端开发，生成前端本地的包，给前端进行引用
