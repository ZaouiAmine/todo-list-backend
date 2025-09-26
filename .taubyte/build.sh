#!/bin/bash

go mod tidy
go build -o /out/function.wasm .
exit 0
