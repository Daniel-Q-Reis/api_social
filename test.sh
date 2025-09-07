#!/bin/bash

# Script to run all tests

echo "Running unit tests..."
go test ./internal/usecase/... -v

if [ $? -eq 0 ]; then
    echo "Unit tests passed!"
else
    echo "Unit tests failed!"
    exit 1
fi

echo "Running integration tests..."
cd integration-test && go test -v

if [ $? -eq 0 ]; then
    echo "Integration tests passed!"
else
    echo "Integration tests failed!"
    exit 1
fi