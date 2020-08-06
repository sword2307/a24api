package a24apiclient

import (
    "fmt"
    "encoding/json"
)

// --------------------------------------------------------------------------------------------------------------------
// List domains
// --------------------------------------------------------------------------------------------------------------------

func (c *T_A24ApiClient) DnsListDomains() (int, interface {}, error) {
    rc, rb, err := c.doApiRequest("GET", c.Config["endpoint"] + "/dns/domains/v1", nil);
    if err != nil {
        return rc, nil, err
    }
    var t T_DnsDomainList
    err = json.Unmarshal([]byte(rb), &t)
    if err != nil {
        return rc, nil, err
    }
    return rc, t, err
}

// --------------------------------------------------------------------------------------------------------------------
// List domain records
// --------------------------------------------------------------------------------------------------------------------

func (c *T_A24ApiClient) DnsListRecords(data map[string]string) (int, T_DnsRecordList, error) {
    rc, rb, err := c.doApiRequest("GET", c.Config["endpoint"] + "/dns/" + data["0"] + "/records/v1", nil);
    if err != nil {
        return rc, nil, err
    }
    var t T_DnsRecordList
    err = json.Unmarshal([]byte(rb), &t)
    if err != nil {
        return rc, nil, err
    }
    return rc, t, err
}

// --------------------------------------------------------------------------------------------------------------------
// Delete dns record
// --------------------------------------------------------------------------------------------------------------------

func (c *T_A24ApiClient) DnsDelete(record interface{}) (int, []byte, error) {

    var lDomain string
    var lHashId string

    switch t := record.(type) {
        case T_DnsRecordA:
            lRecord := record.(T_DnsRecordA)
            lDomain = lRecord.Domain
            lHashId = lRecord.HashId
        case map[string]string:
            lRecord := record.(map[string]string)
            lDomain = lRecord["Domain"]
            lHashId = lRecord["HashId"]
        default:
            return 0, nil, NewA24ApiClientError(fmt.Sprintf("Error: Unknown type %s.", t))
    }

    rc, rb, err :=  c.doApiRequest("DELETE", c.Config["endpoint"] + "/dns/" + lDomain + "/" + lHashId + "/v1", nil);
    if err != nil {
        return rc, nil, err
    }
    return rc, rb, nil
}
