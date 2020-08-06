package a24apiclient

import (
)

type T_DnsDomainList []string

type T_DnsRecordList []map[string]interface{}

type T_DnsRecordCNAME struct {
    Domain         string
    HashId         string     `json:"hashId"`
    Type           string     `json:"type"`
    Name           string     `json:"name"`
    Alias          string     `json:"alias"`
    Ttl            float64    `json:"ttl"`
}

type T_DnsRecordTXT struct {
    Domain         string
    HashId         string     `json:"hashId"`
    Type           string     `json:"type"`
    Name           string     `json:"name"`
    Text           string     `json:"text"`
    Ttl            float64    `json:"ttl"`
}

type T_DnsRecordNS struct {
    Domain         string
    HashId         string     `json:"hashId"`
    Type           string     `json:"type"`
    Name           string     `json:"name"`
    NameServer     string     `json:"nameServer"`
    Ttl            float64    `json:"ttl"`
}
