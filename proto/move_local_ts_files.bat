@echo off
setlocal EnableDelayedExpansion

:: 处理 ts-sea-proto 目录
if exist ".\local-gen\ts-sea-proto" (
    for /f "delims=" %%f in ('dir /b /s ".\local-gen\ts-sea-proto\*.ts" 2^>nul') do (
        set "filePath=%%f"
        set "fileName=%%~nxf"
        set "modulePath=%%~dpf"
        set "modulePath=!modulePath:*definitions\=!"
        
        if not exist ".\local-gen\ts\!modulePath!" (
            mkdir ".\local-gen\ts\!modulePath!" 2>nul
        )
        
        copy /Y "%%f" ".\local-gen\ts\!modulePath!" >nul || (echo [Error] Failed to copy TypeScript file %%f && exit /b 1)
    )
    
    rd /s /q ".\local-gen\ts-sea-proto" 2>nul
)

:: 处理 definitions 目录
if exist ".\local-gen\ts\definitions" (
    for /f "delims=" %%f in ('dir /b /s ".\local-gen\ts\definitions\*.ts" 2^>nul') do (
        set "fullPath=%%f"
        set "relativePath=!fullPath:*\definitions\=!"
        set "targetDir=.\local-gen\ts\!relativePath:~0,-3!"
        
        :: 提取目标目录
        for %%d in ("!targetDir!") do set "targetDir=%%~dpd"
        set "targetDir=!targetDir:~0,-1!"
        
        if not exist "!targetDir!" (
            mkdir "!targetDir!" 2>nul
        )
        
        copy /Y "!fullPath!" "!targetDir!" >nul || (echo [Error] Failed to copy TypeScript file !fullPath! && exit /b 1)
    )
    
    rd /s /q ".\local-gen\ts\definitions" 2>nul
)

:: 删除不需要的目录
if exist ".\local-gen\ts\google" rd /s /q ".\local-gen\ts\google" 2>nul
if exist ".\local-gen\ts\readpaper" rd /s /q ".\local-gen\ts\readpaper" 2>nul

exit /b 0
