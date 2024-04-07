#!/usr/bin/sh
docker compose down
docker compose up --detach
export PATH="$PATH:$HOME/.go/bin"
GOOS=js GOARCH=wasm go build -o web/app.wasm
go run .
