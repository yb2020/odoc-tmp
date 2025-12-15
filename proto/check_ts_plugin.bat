@echo off
echo [Checking] TypeScript plugin
if exist "D:\nodejs\protoc-gen-ts.cmd" (
    echo [Success] TypeScript plugin found at D:\nodejs\protoc-gen-ts.cmd
    exit /b 0
) else (
    echo [Error] protoc-gen-ts not found, please install it first
    echo [Install] npm install --global @protobuf-ts/plugin
    exit /b 1
)
