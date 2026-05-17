#!/bin/bash

echo "========================================"
echo "  Tool Go - Starting Development"
echo "========================================"
echo

echo "[1/3] Checking Go installation..."
go version
if [ $? -ne 0 ]; then
    echo "ERROR: Go is not installed or not in PATH"
    exit 1
fi

echo
echo "[2/3] Starting backend server..."
go run main.go &
BACKEND_PID=$!

echo
echo "[3/3] Starting frontend dev server..."
cd web && npm run dev &
FRONTEND_PID=$!

echo
echo "========================================"
echo "  Backend: http://127.0.0.1:8000"
echo "  Frontend: http://127.0.0.1:3000"
echo "  Swagger: http://127.0.0.1:8000/swagger"
echo "========================================"
echo

trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit" INT TERM

wait
