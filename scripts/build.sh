# /bin/bash

GOOS=js GOARCH=wasm go build -o dist/main.wasm main.go
yarn build

