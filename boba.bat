@echo off
if not exist boba-text.exe (
    echo Building boba-text...
    go build -o boba-text.exe main.go
)
start /b /wait boba-text.exe %*
