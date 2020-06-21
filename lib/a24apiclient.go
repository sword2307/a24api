package a24apiclient

import (
    "os"
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net"
    "time"
    "net/http"
    "path/filepath"
    "strconv"
    "text/tabwriter"
    "regexp"
    "strings"
)

// =============================================================================================================================================================
// CONST
// =============================================================================================================================================================

var (
    C_a24apiclient_config = map[string]string {
        "endpoint": "https://sandboxapi.active24.com",
        "token": "123456qwerty-ok",
        "dial": "tcp",
        "timeout": "30",
    }
    C_a24apiclient_codes = map[string]map[string]map[int]string {
        "_shared_": map[string]map[int]string {
            "_codes_": map[int]string {
                200: "OK",
                204: "OK",
                401: "TOKEN_INVALID",
                403: "UNAUTHORIZED",
                429: "TOO_MANY_REQUESTS",
                500: "SYSTEM_ERROR",
            },
        },
        "dns": map[string]map[int]string {
            "delete": map[int]string {
                400: "DNS_RECORD_TO_DELETE_NOT_FOUND",
            },
            "update": map[int]string {
                400: "DNS_RECORD_TO_UPDATE_NOT_FOUND",
            },
            "create": map[int]string {
                400: "VALIDATION_ERROR",
            },
        },
        "domains": map[string]map[int]string {
            "detail": map[int]string {
                400: "OBJECT_ID_DOESNT_EXIST",
            },
        },
    }
)

// =============================================================================================================================================================
// TYPES
// =============================================================================================================================================================

type TA24ApiClient struct {
    Config                          map[string]string
}

func NewA24ApiClient(config map[string]string) *TA24Client {
    c := &TA24ApiClient{ Config: config }
    c.mergeConfig()

    return c
}

func (c *TA24ApiClient) mergeConfig() {
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

func newA24HttpTransport(timeout, dial string) *http.Transport {
    t := &http.Transport{
        Dial: (func(network, addr string) (net.Conn, error) {
            return (&net.Dialer{
                Timeout:        strconv.Atoi(c.Config["timeout"]) * time.Second,
                LocalAddr:      nil,
                DualStack:      false,
            }).Dial(c.Config["dial"], addr)
        }),
    }
    return t
}

func newA24HttpClient(http_transport *http.Transport) *http.Client {
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

    if a24api_request, err := http.NewRequest("GET", c.Config["endpoint"] + "/dns/domains/v1", bytes.NewBuffer(a24api_request_body_json)); err != nil {
        return (nil, err)
    }

    a24api_request.Header.Set("Content-type", "application/json")
    a24api_request.Header.Set("Accept", "application/json")
    a24api_request.Header.Set("Authorization", "Bearer " + c.Config["token"])

    if a24api_response, err := a24api_client.Do(a24api_request); err != nil {
        return (nil, error)
    }

    defer a24api_response.Body.Close()

    if a24api_response_body, err := ioutil.ReadAll(a24api_response.Body); err != nil {
        return (nil, error)
    } else {
        return (a24api_response_body, nil)
    }
}

// List dns records
func (c *TA24ApiClient) DnsRecords(domain string) ([]byte, error) {

}

// Create dns record
func (c *TA24ApiClient) DnsCreate(record interface) ([]byte, error) {

}

// Update dns record
func (c *TA24ApiClient) DnsUpdate(record interface) ([]byte, error) {

}

// Delete dns record
func (c *TA24ApiClient) DnsDelete(record interface) ([]byte, error) {

}

// ================================================================================================================================================================
// LOAD CONFIG FILE
// ================================================================================================================================================================

    configFile, err := os.Open(a24api["config"])
    if err == nil {
        defer configFile.Close()
        configDataRaw, _ := ioutil.ReadAll(configFile)
        var configData map[string]string
        json.Unmarshal(configDataRaw, &configData)
        if configData["a24api_endpoint"] != "" {
            a24api["endpoint"] = configData["a24api_endpoint"]
        }
        if configData["a24api_token"] != "" {
            a24api["token"] = configData["a24api_token"]
        }
    }

// ================================================================================================================================================================
// CHECK INPUT DATA
// ================================================================================================================================================================

    if a24api["service"] == "" || a24api["function"] == "" {
        fmt.Println("Service or function not provided.")
        printHelp()
        os.Exit(1)
    }

// ================================================================================================================================================================
// PREPARE REQUEST
// ================================================================================================================================================================

    a24api_request_body := make(map[string]string)
    var posArgOffset = 0

    switch a24api["service"] {
        case "dns":
            switch a24api["function"] {
                case "list":
                    if len(a24api_args) == 0 {
                        a24api["endpoint-uri"] = "/dns/domains/v1"
                        a24api["endpoint-method"] = "GET"
                    // expected arguments: 0=domain
                    } else {
                        a24api["domain"] = a24api_args["argument0"]
                        a24api["endpoint-uri"] = "/dns/" + a24api["domain"] + "/records/v1"
                        a24api["endpoint-method"] = "GET"
                    }
                case "create", "update":
                    if a24api["function"] == "create" {
                        a24api["endpoint-method"] = "POST"
                    // expected arguments: 0=domain, 1=hash_id
                    } else {
                        a24api["endpoint-method"] = "PUT"
                        a24api_request_body["hashId"] = a24api_args["argument1"]
                        posArgOffset = 1
                    }
                    // expected arguments: 0=domain, 1(2)=type
                    a24api["domain"] = a24api_args["argument0"]
                    a24api["record-type"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 1)]
                    a24api["endpoint-uri"] = "/dns/" + a24api["domain"] + "/" + strings.ToLower(a24api["record-type"]) + "/v1"
                    switch a24api["record-type"] {
                        case "A", "AAAA":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=value
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["ip"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                        case "CNAME":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=value
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["alias"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                        case "TXT":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=value
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["text"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                        case "NS":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=value
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["nameServer"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                        case "SSHFP":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=algorithm, 5(6)=fp_type, 6(7)=fingerprint
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["algorithm"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                            a24api_request_body["fingerprintType"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 5)]
                            a24api_request_body["text"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 6)]
                        case "SRV":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=priority, 5(6)=weight, 6(7)=port, 7(8)=target
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["priority"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                            a24api_request_body["weight"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 5)]
                            a24api_request_body["port"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 6)]
                            a24api_request_body["target"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 7)]
                        case "TLSA":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=certificate_usage, 5(6)=selector, 6(7)=matching_type, 7(8)=hash
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["certificateUsage"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                            a24api_request_body["selector"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 5)]
                            a24api_request_body["matchingType"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 6)]
                            a24api_request_body["hash"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 7)]
                        case "CAA":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=flags, 5(6)=tag, 6(7)=value
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["flags"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                            a24api_request_body["tag"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 5)]
                            a24api_request_body["caaValue"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 6)]
                        case "MX":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=priority, 5(6)=value
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["priority"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                            a24api_request_body["mailserver"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 5)]
                        default:
                            fmt.Printf("Unsupported dns type: %s.\n", a24api["record-type"])
                            os.Exit(1)
                    }
                // expected arguments: 0=domain, 1=hash_id
                case "delete":
                    a24api["domain"] = a24api_args["argument0"]
                    a24api["endpoint-uri"] = "/dns/" + a24api["domain"] + "/" + a24api_args["argument1"] + "/v1"
                    a24api["endpoint-method"] = "DELETE"
                    a24api_request_body["hashId"] = a24api_args["argument1"]
                default:
                    fmt.Printf("Unsupported function: %s.\n", a24api["function"])
                    os.Exit(1)
            }
        default:
            fmt.Printf("Unsupported service: %s.", a24api["service"])
            os.Exit(1)
    }

// ================================================================================================================================================================
// MAKE REQUEST
// ================================================================================================================================================================

    a24api_transport := &http.Transport{
        Dial: (func(network, addr string) (net.Conn, error) {
            return (&net.Dialer{
                Timeout:        10 * time.Second,
                LocalAddr:      nil,
                DualStack:      false,
            }).Dial(a24api["dial"], addr)
        }),
    }

    a24api_client := &http.Client{Transport: a24api_transport}

    a24api_request_body_json, err := json.Marshal(a24api_request_body)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    a24api_request, err := http.NewRequest(a24api["endpoint-method"], a24api["endpoint"] + a24api["endpoint-uri"], bytes.NewBuffer(a24api_request_body_json))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    a24api_request.Header.Set("Content-type", "application/json")
    a24api_request.Header.Set("Accept", "application/json")
    a24api_request.Header.Set("Authorization", "Bearer " + a24api["token"])

    a24api_response, err := a24api_client.Do(a24api_request)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    defer a24api_response.Body.Close()

    a24api_response_body, err := ioutil.ReadAll(a24api_response.Body)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
