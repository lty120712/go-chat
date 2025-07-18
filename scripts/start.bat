@echo off
REM ========================================
REM 启动 Go 项目
REM ========================================

REM 设置项目根目录
SET PROJECT_DIR=%cd%

REM 设置命令行使用 UTF-8 编码，避免中文乱码
chcp 65001

REM 显示当前目录路径
echo 当前工作目录: %PROJECT_DIR%

REM 进入 cmd 目录
cd ..\cmd
IF %ERRORLEVEL% NEQ 0 (
    echo 进入 cmd 目录失败！
    exit /b %ERRORLEVEL%
)

REM 使用 fresh 启动项目（热重载）
echo 启动项目...
fresh

REM 检查 fresh 是否成功启动
IF %ERRORLEVEL% NEQ 0 (
    echo 项目启动失败，请检查错误日志。
    exit /b %ERRORLEVEL%
)

REM 启动完成，暂停等待查看日志
echo 项目已启动。按任意键退出...
pause
