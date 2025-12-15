@echo off
REM Step 1: Create package-lock.json
echo [INFO] Creating package-lock.json...
(
    echo {
    echo   "name": "go-sea-proto",
    echo   "version": "0.0.70",
    echo   "lockfileVersion": 3,
    echo   "requires": true,
    echo   "packages": {
    echo     "": {
    echo       "name": "go-sea-proto",
    echo       "version": "0.0.70",
    echo       "dependencies": {
    echo         "@protobuf-ts/runtime": "^2.11.0"
    echo       }
    echo     },
    echo     "node_modules/@protobuf-ts/runtime": {
    echo       "version": "2.11.0",
    echo       "resolved": "https://registry.npmmirror.com/@protobuf-ts/runtime/-/runtime-2.11.0.tgz",
    echo       "integrity": "sha512-DfpRpUiNvPC3Kj48CmlU4HaIEY1Myh++PIumMmohBAk8/k0d2CkxYxJfPyUAxfuUfl97F4AvuCu1gXmfOG7OJQ==",
    echo       "license": "(Apache-2.0 AND BSD-3-Clause)"
    echo     }
    echo   }
    echo }
) > package-lock.json
echo [SUCCESS] package-lock.json created.

REM Step 2: Run make command
echo.
echo [INFO] Running build command: mingw32-make -f makefileForWindows all
mingw32-make -f makefileForWindows all
if %errorlevel% neq 0 (
    echo [ERROR] Build command failed.
    exit /b 1
)
echo [SUCCESS] Build command finished.

REM Step 3: Run npm commands
echo.
echo [INFO] Running npm link...
call npm link
if %errorlevel% neq 0 (
    echo [ERROR] 'npm link' failed.
    exit /b 1
)
echo [SUCCESS] 'npm link' finished.

echo.
echo [INFO] Running npm install @protobuf-ts/runtime...
call npm install @protobuf-ts/runtime
if %errorlevel% neq 0 (
    echo [ERROR] 'npm install' failed.
    exit /b 1
)
echo [SUCCESS] 'npm install' finished.

echo.
echo All steps completed successfully.
