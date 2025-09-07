@echo off
echo Running unit tests...
go test ./internal/usecase/... -v
if %errorlevel% neq 0 (
    echo Unit tests failed!
    exit /b %errorlevel%
)
echo Unit tests passed!

echo Running integration tests...
cd integration-test
go test -v
if %errorlevel% neq 0 (
    echo Integration tests failed!
    exit /b %errorlevel%
)
echo Integration tests passed!