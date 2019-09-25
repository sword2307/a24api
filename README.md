# a24api
Active24 REST API Client Implementation In GO



#### Build targets:
linux-386
linux-amd64
linux-arm
linux-arm64
windows-386
windows-amd64
darwin-386
darwin-amd64

#### Build tips:
disable symbol table and DWARF generation

`go build -ldflags="-s -w"`


compress binary

`/usr/bin/upx --brute a24api`
