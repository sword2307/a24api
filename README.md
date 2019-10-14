# a24api
Active24 REST API Command Line Client Implementation In GO

API Provider: ACTIVE 24, s.r.o

API URL: https://sandboxapi.active24.com


#### Implemented services/functions
- dns
    - list
    - create A,AAAA,CNAME,TXT,NS,SSHFP,SRV,TLSA,CAA,MX
    - update A,AAAA,CNAME,TXT,NS,SSHFP,SRV,TLSA,CAA,MX
    - delete


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

`CGO_ENABLED=0 go build -ldflags="-s -w"`

compress binary

`/usr/bin/upx --brute a24api`


#### Docker

docker build --tag a24api .

docker run --rm a24api
