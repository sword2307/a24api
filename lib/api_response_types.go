package a24apiclient

import (
)

type T_DnsRecordA struct {
    HashId         string
    Type           string
    Ip             string
    Name           string
    Ttl            float64
}

type T_DnsRecordAAAA struct {
    hashId         string
    Type           string
    ip             string
    name           string
    ttl            float64
}

type T_DnsRecordCNAME struct {
    HashId         string
    Type           string
    Name           string
    Alias          string
    Ttl            float64
}

type T_DnsRecordTXT struct {
    HashId         string
    Type           string
    Name           string
    Text           string
    Ttl            float64
}

type T_DnsRecordNS struct {
    HashId         string
    Type           string
    Name           string
    NameServer     string
    Ttl            float64
}
