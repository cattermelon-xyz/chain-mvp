#!/bin/bash
go build -o ./bin/htg-client build/client/main.go
go build -o ./bin/htg build/server/main.go
go build -o ./bin/demo build/main.go