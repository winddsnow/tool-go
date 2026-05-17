@echo off
echo ========================================
echo   Tool Go - Starting Development
echo ========================================
echo.

echo [1/3] Checking Go installation...
go version
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed or not in PATH
    pause
    exit /b 1
)

echo.
echo [2/3] Starting backend server...
start cmd /k "cd /d %~dp0 && go run main.go"

echo.
echo [3/3] Starting frontend dev server...
cd web
start cmd /k "npm run dev"

echo.
echo ========================================
echo   Backend: http://127.0.0.1:8000
echo   Frontend: http://127.0.0.1:3000
echo   Swagger: http://127.0.0.1:8000/swagger
echo ========================================
echo.
pause
