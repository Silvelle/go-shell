@echo off
rem Builds and starts the emulator on Windows without make.
rem Needs Go and a C compiler (gcc) in PATH, because the GUI library uses C.
rem Arguments go to the emulator, for example: run.bat --vfs vfs\deep
set CGO_ENABLED=1
if not exist bin mkdir bin
go build -o bin\emulator.exe .\src || exit /b 1
bin\emulator.exe %*
