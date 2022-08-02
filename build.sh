#!/bin/sh

# https://github.com/xaionaro/documentation/blob/master/golang/reduce-binary-size.md
# go build -ldflags="-w -s" -gcflags=all="-l -B" .
# upx --best --ultra-brute gotem
# upx --best --lzma gotem
go build -ldflags="-s -w" .
upx gotem
