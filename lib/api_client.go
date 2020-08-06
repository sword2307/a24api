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
    C_A24ApiClient_Config = map[string]string {
        "endpoint": "https://sandboxapi.active24.com",
        "token": "123456qwerty-ok",
        "network": "tcp",                                // [tcp|tcp4|tcp6]
        "timeout": "30",
    }
)

// =============================================================================================================================================================
// TYPES
// =============================================================================================================================================================

type T_A24ApiClient_Config_File struct {
    Endpoint   string
    Token      string
    Network    string
    Timeout    string
}

type T_A24ApiClient struct {
    Config                          map[string]string
    HttpClient                      *http.Client
}

type T_A24ApiClientError struct{
    Text string
}

func (e *T_A24ApiClientError) Error() string {
    return e.Text
}

func NewA24ApiClientError(text string) *T_A24ApiClientError {
    e:= &T_A24ApiClientError{ Text: text }
    return e
}

func NewA24ApiClient(config map[string]string) *T_A24ApiClient {
    c := &T_A24ApiClient{ Config: config }
    c.mergeConfig()
    c.HttpClient = c.newHttpClient()

    return c
}

func (c *T_A24ApiClient) mergeConfig() {
    if c.Config == nil {
        c.Config = make(map[string]string)
    }
    for key, defaultValue := range C_A24ApiClient_Config {
        if _, isPresent := c.Config[key]; !isPresent {
            c.Config[key] = defaultValue
        }
    }
}

func (c *T_A24ApiClient) GetCodeText(code int, service string, function string) (code_text string) {
    if code_text = C_A24ApiClient_Codes[service][function][code]; code_text == "" {
        if code_text = C_A24ApiClient_Codes["_shared_"]["_codes_"][code]; code_text == "" {
            code_text = "UNKNOWN_CODE"
        }
    }
    return
}

func newHttpTransport(timeout, network string) *http.Transport {
    var l_dualstack bool = true
    l_timeout_i, _ := strconv.Atoi(timeout)
    l_timeout_d := time.Duration(l_timeout_i) * time.Second
    if !(network == "tcp") {
        l_dualstack = false
    }
    t := &http.Transport{
        Dial: (func(network, addr string) (net.Conn, error) {
            return (&net.Dialer{
                Timeout:        l_timeout_d,
                LocalAddr:      nil,
                DualStack:      l_dualstack,
            }).Dial(network, addr)
        }),
    }
    return t
}

func (c *T_A24ApiClient) newHttpClient() *http.Client {
    http_transport := newHttpTransport(c.Config["network"], c.Config["timeout"])
    s := &http.Client{Transport: http_transport}
    return s
}

// =============================================================================================================================================================
// API FUNCTIONS
// =============================================================================================================================================================

func (c *T_A24ApiClient) doApiRequest(method, endpoint string, body map[string]string) (int, []byte, error) {

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
