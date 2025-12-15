@echo off
setlocal EnableDelayedExpansion

:: 处理 definitions 目录
if exist ".\gen\ts\definitions" (
    for /f "delims=" %%f in ('dir /b /s ".\gen\ts\definitions\*.ts" 2^>nul') do (
        set "fullPath=%%f"
        set "relativePath=!fullPath:*\definitions\=!"
        set "targetDir=.\gen\ts\!relativePath:~0,-3!"
        
        :: 提取目标目录
        for %%d in ("!targetDir!") do set "targetDir=%%~dpd"
        set "targetDir=!targetDir:~0,-1!"
        
        if not exist "!targetDir!" (
            mkdir "!targetDir!" 2>nul
        )
        
        copy /Y "!fullPath!" "!targetDir!" >nul || (echo [Error] Failed to copy TypeScript file !fullPath! && exit /b 1)
    )
    
    rd /s /q ".\gen\ts\definitions" 2>nul
)

:: 删除不需要的目录
if exist ".\gen\ts\google" rd /s /q ".\gen\ts\google" 2>nul
if exist ".\gen\ts\readpaper" rd /s /q ".\gen\ts\readpaper" 2>nul

exit /b 0
