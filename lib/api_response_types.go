package a24apiclient

import (
)

type T_DnsDomainList []string

type T_DnsRecordList []map[string]interface{}

type T_DnsRecordA struct {
    HashId          string    `json:"hashId"`
    Type            string    `json:"type"`
    Ip              string    `json:"ip"`
    Name            string    `json:"name"`
    Ttl             float64   `json:"ttl"`
}

type T_DnsRecordAAAA struct {
    hashId         string     `json:"hashId"`
    Type           string     `json:"type"`
    ip             string     `json:"ip"`
    name           string     `json:"name"`
    ttl            float64    `json:"ttl"`
}

type T_DnsRecordCNAME struct {
    HashId         string     `json:"hashId"`
    Type           string     `json:"type"`
    Name           string     `json:"name"`
    Alias          string     `json:"alias"`
    Ttl            float64    `json:"ttl"`
}

type T_DnsRecordTXT struct {
    HashId         string     `json:"hashId"`
    Type           string     `json:"type"`
    Name           string     `json:"name"`
    Text           string     `json:"text"`
    Ttl            float64    `json:"ttl"`
}

type T_DnsRecordNS struct {
    HashId         string     `json:"hashId"`
    Type           string     `json:"type"`
    Name           string     `json:"name"`
    NameServer     string     `json:"nameServer"`
    Ttl            float64    `json:"ttl"`
}
