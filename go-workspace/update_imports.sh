#!/bin/bash

# Update imports in all Go files
find /home/rohankarmacharya/Task\ 1/go-workspace/movie-lib -type f -name "*.go" -exec sed -i 's|github.com/rohankarmacharya/movie-lib|github.com/rohankarmacharya/go-workspace/movie-lib|g' {} \;

# Run go mod tidy to update dependencies
cd "/home/rohankarmacharya/Task 1/go-workspace/movie-lib" && go mod tidy
