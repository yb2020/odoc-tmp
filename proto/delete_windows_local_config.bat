@echo off
echo [INFO] Starting cleanup process...

REM Delete gen directory
if exist "gen" (
    echo [INFO] Removing gen directory...
    rd /s /q gen
    if %errorlevel% neq 0 (
        echo [ERROR] Failed to remove gen directory.
        exit /b 1
    )
    echo [SUCCESS] gen directory removed.
) else (
    echo [INFO] gen directory not found, skipping.
)

REM Delete node_modules directory
if exist "node_modules" (
    echo [INFO] Removing node_modules directory...
    rd /s /q node_modules
    if %errorlevel% neq 0 (
        echo [ERROR] Failed to remove node_modules directory.
        exit /b 1
    )
    echo [SUCCESS] node_modules directory removed.
) else (
    echo [INFO] node_modules directory not found, skipping.
)

REM Delete package-lock.json file
if exist "package-lock.json" (
    echo [INFO] Removing package-lock.json file...
    del /f package-lock.json
    if %errorlevel% neq 0 (
        echo [ERROR] Failed to remove package-lock.json file.
        exit /b 1
    )
    echo [SUCCESS] package-lock.json file removed.
) else (
    echo [INFO] package-lock.json file not found, skipping.
)

echo.
echo [INFO] Cleanup completed successfully.
