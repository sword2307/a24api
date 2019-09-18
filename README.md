# a24api
Active24 REST API Client Implementation In GO

#### Build tips:
disable symbol table and DWARF generation
`go build -ldflags="-s -w"`
compress binary
`/usr/bin/upx --brute a24api`
