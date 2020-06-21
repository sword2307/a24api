package a24apiclient

type DnsRecordA struct {
    hashId         string
    type           string
    ip             string
    name           string
    ttl            float64
}

type DnsRecordAAAA struct {
    hashId         string
    type           string
    ip             string
    name           string
    ttl            float64
}

type DnsRecordCNAME struct {
    hashId         string
    type           string
    name           string
    alias          string
    ttl            float64
}

type DnsRecordTXT struct {
    hashId         string
    type           string
    name           string
    text           string
    ttl            float64
}

type DnsRecordNS struct {
    hashId         string
    type           string
    name           string
    nameServer     string
    ttl            float64
}
