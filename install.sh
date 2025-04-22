#!/bin/bash
echo "go build -o $(go env GOPATH)/bin/mcp-moedas cmd/mcp/main.go"
go build -o $(go env GOPATH)/bin/mcp-moedas cmd/mcp/main.go