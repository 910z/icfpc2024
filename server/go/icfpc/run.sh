#!/usr/bin/sh
GOOS=js GOARCH=wasm go build -o web/app.wasm
go run .
