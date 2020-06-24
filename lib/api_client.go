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

type T_a24apiclient struct {
    Config                          map[string]string
    HttpClient                      *http.Client
}

func New_a24apiclient(config map[string]string) *T_a24apiclient {
    c := &T_a24apiclient{ Config: config }
    c.mergeConfig()
    c.HttpClient = c.newHttpClient()

    return c
}

func (c *T_a24apiclient) mergeConfig() {
    if c.Config == nil {
        c.Config = make(map[string]string)
    }
    for key, defaultValue := range C_a24apiclient_config {
        if _, isPresent := c.Config[key]; !isPresent {
            c.Config[key] = defaultValue
        }
    }
}

func (c *T_a24apiclient) getCodeText(code int, service string, function string) (code_text string) {
    if code_text = C_a24apiclient_codes[service][function][code]; code_text == "" {
        if code_text = C_a24apiclient_codes["_shared_"]["_codes_"][code]; code_text == "" {
            code_text = "UNKNOWN_CODE"
        }
    }
    return
}

func newHttpTransport(timeout, network string) *http.Transport {
    l_timeout_i, _ := strconv.Atoi(timeout)
    l_timeout_d := time.Duration(l_timeout_i) * time.Second
    t := &http.Transport{
        Dial: (func(network, addr string) (net.Conn, error) {
            return (&net.Dialer{
                Timeout:        l_timeout_d,
                LocalAddr:      nil,
                DualStack:      true,
            }).Dial(network, addr)
        }),
    }
    return t
}

func (c *T_a24apiclient) newHttpClient() *http.Client {
    http_transport := newHttpTransport(C_a24apiclient_config["network"], C_a24apiclient_config["timeout"])
    s := &http.Client{Transport: http_transport}
    return s
}

// =============================================================================================================================================================
// API FUNCTIONS
// =============================================================================================================================================================

func (c *T_a24apiclient) doApiRequest(method, endpoint string, body map[string]string) (int, []byte, error) {

    body_json, err := json.Marshal(body)
    if err != nil {
        return 0, nil, err
    }

    a24api_request, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body_json))
    if err != nil {
        return 0, nil, err
    }

    a24api_request.Header.Set("Content-type", "application/json")
    a24api_request.Header.Set("Accept", "application/json")
    a24api_request.Header.Set("Authorization", "Bearer " + c.Config["token"])

    a24api_response, err := c.HttpClient.Do(a24api_request)
    if err != nil {
        return 0, nil, err
    }

    defer a24api_response.Body.Close()

    if a24api_response_body, err := ioutil.ReadAll(a24api_response.Body); err != nil {
        return 0, nil, err
    } else {
        return a24api_response.StatusCode, a24api_response_body, nil
    }

}

// List domains
func (c *T_a24apiclient) DnsListDomains() (int, []byte, error) {

    a24api_response_code, a24api_response_body, err := c.doApiRequest("GET", c.Config["endpoint"] + "/dns/domains/v1", nil);

    return a24api_response_code, a24api_response_body, err

}

// List domain dns records
func (c *T_a24apiclient) DnsListRecords(domain string) (int, []byte, error) {

    a24api_response_code, a24api_response_body, err := c.doApiRequest("GET", c.Config["endpoint"] + "/dns/" + domain + "/records/v1", nil);

    return a24api_response_code, a24api_response_body, err

}

// Create dns record
func (c *T_a24apiclient) DnsCreate(record interface {}) (int, []byte, error) {

    return 0, nil, nil

}

// Update dns record
func (c *T_a24apiclient) DnsUpdate(record interface {}) (int, []byte, error) {

    return 0, nil, nil

}

// Delete dns record
func (c *T_a24apiclient) DnsDelete(record interface {}) (int, []byte, error) {

    return 0, nil, nil

}
