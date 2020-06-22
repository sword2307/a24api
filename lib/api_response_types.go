package a24apiclient

import (
)

type DnsRecordA struct {
    HashId         string
    Type           string
    Ip             string
    Name           string
    Ttl            float64
}

type DnsRecordAAAA struct {
    hashId         string
    Type           string
    ip             string
    name           string
    ttl            float64
}

type DnsRecordCNAME struct {
    HashId         string
    Type           string
    Name           string
    Alias          string
    Ttl            float64
}

type DnsRecordTXT struct {
    HashId         string
    Type           string
    Name           string
    Text           string
    Ttl            float64
}

type DnsRecordNS struct {
    HashId         string
    Type           string
    Name           string
    NameServer     string
    Ttl            float64
}
