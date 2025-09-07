@echo off
echo Running database integration tests...
go test -v -run "TestPostgres"
if %errorlevel% neq 0 (
    echo Database integration tests failed!
    exit /b %errorlevel%
)
echo Database integration tests passed!