@echo off
setlocal EnableDelayedExpansion

:: 处理 definitions 目录
if exist ".\local-gen\go\definitions" (
    for /f "delims=" %%f in ('dir /b /s ".\local-gen\go\definitions\*.pb.go" 2^>nul') do (
        set "fullPath=%%f"
        set "relativePath=!fullPath:*\definitions\=!"
        set "targetDir=.\local-gen\go\!relativePath:~0,-7!"
        
        :: 提取目标目录
        for %%d in ("!targetDir!") do set "targetDir=%%~dpd"
        set "targetDir=!targetDir:~0,-1!"
        
        if not exist "!targetDir!" (
            mkdir "!targetDir!" 2>nul
        )
        
        copy /Y "!fullPath!" "!targetDir!" >nul || (echo [Error] Failed to copy Go file !fullPath! && exit /b 1)
    )
    
    :: 处理验证文件
    for /f "delims=" %%f in ('dir /b /s ".\local-gen\go\definitions\*.validate.pb.go" ".\local-gen\go\definitions\*_validate.pb.go" 2^>nul') do (
        set "fullPath=%%f"
        set "relativePath=!fullPath:*\definitions\=!"
        set "targetDir=.\local-gen\go\!relativePath:~0,-15!"
        
        :: 提取目标目录
        for %%d in ("!targetDir!") do set "targetDir=%%~dpd"
        set "targetDir=!targetDir:~0,-1!"
        
        if not exist "!targetDir!" (
            mkdir "!targetDir!" 2>nul
        )
        
        copy /Y "!fullPath!" "!targetDir!" >nul || (echo [Error] Failed to copy Go validation file !fullPath! && exit /b 1)
    )
    
    :: 处理 gRPC 文件
    for /f "delims=" %%f in ('dir /b /s ".\local-gen\go\definitions\*_grpc.pb.go" 2^>nul') do (
        set "fullPath=%%f"
        set "relativePath=!fullPath:*\definitions\=!"
        set "targetDir=.\local-gen\go\!relativePath:~0,-12!"
        
        :: 提取目标目录
        for %%d in ("!targetDir!") do set "targetDir=%%~dpd"
        set "targetDir=!targetDir:~0,-1!"
        
        if not exist "!targetDir!" (
            mkdir "!targetDir!" 2>nul
        )
        
        copy /Y "!fullPath!" "!targetDir!" >nul || (echo [Error] Failed to copy Go gRPC file !fullPath! && exit /b 1)
    )
    
    rd /s /q ".\local-gen\go\definitions" 2>nul
)

:: 删除不需要的目录
if exist ".\local-gen\go\google" rd /s /q ".\local-gen\go\google" 2>nul
if exist ".\local-gen\go\readpaper" rd /s /q ".\local-gen\go\readpaper" 2>nul

exit /b 0
