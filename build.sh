#!/bin/sh

go build -ldflags="-s -w" .
upx gotem
# go build -ldflags="-w -s" -gcflags=all="-l -B" .
# upx --best --ultra-brute gotem
# upx --best --lzma gotem
