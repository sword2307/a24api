package a24apiclient

import (
    "fmt"
    "strconv"
)

// --------------------------------------------------------------------------------------------------------------------
// Type
// --------------------------------------------------------------------------------------------------------------------

type T_DnsRecordA struct {
    Domain          string
    HashId          string    `json:"hashId"`
    Type            string    `json:"type"`
    Ip              string    `json:"ip"`
    Name            string    `json:"name"`
    Ttl             float64   `json:"ttl"`
}

func NewDnsRecordA(data map[string]string) (*T_DnsRecordA, error) {
    r := &T_DnsRecordA{}
    r.Domain = data["Domain"]
    r.HashId = data["HashId"]
    r.Type = data["Type"]
    r.Ip = data["Ip"]
    r.Name = data["Name"]
    lTtl, err := strconv.ParseFloat(data["Ttl"], 64)
    if err != nil {
        return nil, err
    }
    r.Ttl = lTtl
    return r, nil
}

func (c *T_A24ApiClient) DnsCreateUpdateA(record interface{}, action string) (int, []byte, error) {

    var lApiData map[string]string
    lApiData = make(map[string]string)
    var lDomain string
    var lMethod string

    switch t := record.(type) {
        case T_DnsRecordA:
            lRecord := record.(T_DnsRecordA)
            lDomain = lRecord.Domain
            lApiData["name"] = lRecord.Name
            lApiData["ttl"] = fmt.Sprintf("%g", lRecord.Ttl)
            lApiData["ip"] = lRecord.Ip
            if action == "update" {
                lApiData["hashId"] = lRecord.HashId
            }
        case map[string]string:
            lRecord := record.(map[string]string)
            lDomain = lRecord["Domain"]
            lApiData["name"] = lRecord["Name"]
            lApiData["ttl"] = lRecord["Ttl"]
            lApiData["ip"] = lRecord["Ip"]
            if action == "update" {
                lApiData["hashId"] = lRecord["HashId"]
            }
        default:
            return 0, nil, NewA24ApiClientError(fmt.Sprintf("Error: Unknown type %s.", t))
    }

    if action == "update" {
        lMethod = "PUT"
    } else {
        lMethod = "POST"
    }

    rc, rb, err :=  c.doApiRequest(lMethod, c.Config["endpoint"] + "/dns/" + lDomain + "/a/v1", lApiData);
    if err != nil {
        return rc, nil, err
    }
    return rc, rb, nil
}
