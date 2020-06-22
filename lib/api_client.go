package a24apiclient

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net"
    "time"
    "net/http"
    "strconv"
)

// =============================================================================================================================================================
// CONST
// =============================================================================================================================================================

var (
    C_a24apiclient_config = map[string]string {
        "endpoint": "https://sandboxapi.active24.com",
        "token": "123456qwerty-ok",
        "network": "tcp",
        "dualstack": "true",
        "timeout": "30",
    }
)

// =============================================================================================================================================================
// TYPES
// =============================================================================================================================================================

type TA24ApiClient struct {
    Config                          map[string]string
    HttpClient                      *http.Client
}

func NewA24ApiClient(config map[string]string) *TA24Client {
    c := &TA24ApiClient{ Config: config }
    c.mergeConfig()
    c.HttpClient = newHttpClient()

    return c
}

func (c *TA24ApiClient) mergeConfig() {
    if c.Config == nil {
        c.Config = make(map[string]string)
    }
    for key, defaultValue := range C_a24apiclient_config {
        if _, isPresent := s.Config[key]; !isPresent {
            c.Config[key] = defaultValue
        }
    }
}

func (c *TA24ApiClient) getCodeText(code int, service string, function string) (code_text string) {
    if code_text = C_a24api_codes[service][function][code]; code_text == "" {
        if code_text = C_a24api_codes["_shared_"]["_codes_"][code]; code_text == "" {
            code_text = "UNKNOWN_CODE"
        }
    }
    return
}

func newHttpTransport(timeout, network string) *http.Transport {
    t := &http.Transport{
        Dial: (func(network, addr string) (net.Conn, error) {
            return (&net.Dialer{
                Timeout:        strconv.Atoi(timeout) * time.Second,
                LocalAddr:      nil,
                DualStack:      true,
            }).Dial(network, addr)
        }),
    }
    return t
}

func (c *TA24ApiClient) newHttpClient() *http.Client {
    http_transport := newHttpTransport(C_a24apiclient_config["network"], C_a24apiclient_config["timeout"])
    s := &http.Client{Transport: http_transport}
    return s
}

// =============================================================================================================================================================
// API FUNCTIONS
// =============================================================================================================================================================

func (c *TA24ApiClient) doApiRequest(method, endpoint string, body map[string]string) ([]byte, error) {

    if body_json, err := json.Marshal(body); err != nil {
        return nil, err
    }

    if a24api_request, err := http.NewRequest(a24api["endpoint-method"], a24api["endpoint"] + a24api["endpoint-uri"], bytes.NewBuffer(a24api_request_body_json)); err != nil {
        return nil, err
    }

    a24api_request.Header.Set("Content-type", "application/json")
    a24api_request.Header.Set("Accept", "application/json")
    a24api_request.Header.Set("Authorization", "Bearer " + a24api["token"])

    if a24api_response, err := a24api_client.Do(a24api_request); err != nil {
        return nil, err
    }

    defer a24api_response.Body.Close()

    if a24api_response_body, err := ioutil.ReadAll(a24api_response.Body); err != nil {
        return nil, err
    } else {
        return a24api_response_body, nil
    }

}

func (c *TA24ApiClient) DnsDomains() ([]byte, error) {

    if a24api_request, err := c.HttpClient.NewRequest("GET", c.Config["endpoint"] + "/dns/domains/v1", bytes.NewBuffer(a24api_request_body_json)); err != nil {
        return nil, err
    }

    a24api_request.Header.Set("Content-type", "application/json")
    a24api_request.Header.Set("Accept", "application/json")
    a24api_request.Header.Set("Authorization", "Bearer " + c.Config["token"])

    if a24api_response, err := a24api_client.Do(a24api_request); err != nil {
        return nil, error
    }

    defer a24api_response.Body.Close()

    if a24api_response_body, err := ioutil.ReadAll(a24api_response.Body); err != nil {
        return nil, error
    } else {
        return a24api_response_body, nil
    }
}

// List dns records
func (c *TA24ApiClient) DnsRecords(domain string) ([]byte, error) {
    return nil, nil
}

// Create dns record
func (c *TA24ApiClient) DnsCreate(record interface {}) ([]byte, error) {
    return nil, nil
}

// Update dns record
func (c *TA24ApiClient) DnsUpdate(record interface {}) ([]byte, error) {
    return nil, nil
}

// Delete dns record
func (c *TA24ApiClient) DnsDelete(record interface {}) ([]byte, error) {
    return nil, nil
}

