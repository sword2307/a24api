package a24api

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

const (
    C_A24API_ENDPOINT = "https://sandboxapi.active24.com"
    C_A24API_TOKEN = "123456qwerty-ok"
    C_A24API_CONFIG = "a24api-conf.json"
    C_A24API_FORMAT = "inline"
    C_A24API_DIAL = "tcp"
    C_A24API_NAME_REGEXP = ".*"
    C_A24API_TYPE_REGEXP = ".*"
    C_A24API_VALUE_REGEXP = ".*"
)

var C_A24API_CODES = map[string]map[string]map[int]string {
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

type T_A24API_Config_File struct {
    a24api_endpoint string
    a24api_token string
}

func getCodeText(code int, service string, function string) (code_text string) {
    if code_text = Con_a24api_codes[service][function][code]; code_text == "" {
        if code_text = Con_a24api_codes["_shared_"]["_codes_"][code]; code_text == "" {
            code_text = "UNKNOWN_CODE"
        }
    }
    return
}

func printHelp() {
    fmt.Println(`Usage: a24api [options] <service> <function> [parameters]

Options:
    -c|--config <path>            Path to config file. Default is a24api-conf.json. Can be also set via env A24API_CONFIG.
    -e|--endpoint <url>           Active24 REST API url. Can be also set via env A24API_ENDPOINT.
    -t|--token <token>            Active24 REST API token. Can be also set via env A24API_TOKEN.
    -f|--format <json|inline>     Output format (default: inline).
    -4                            Use ipv4.
    -6                            Use ipv6.

Services, functions and parameters:
    dns
        list [-fn <name regex filter>]
        list <domain> [-ft <type regex filter>] [-fn <name regex filter>] [-fv <value regex filter>]
        delete <domain> <hash_id>
        create <domain>
            <A|AAAA|CNAME|TXT> <name|@> <ttl> <ip|alias|text>
            <NS> <name|@> <ttl> <nameserver>
            <SSHFP> <name> <ttl> <algorithm> <fp_type> <fingerprint>
            <SRV> <name> <ttl> <priority> <weight> <port> <target>
            <TLSA> <name> <ttl> <certificate_usage> <selector> <matching_type> <hash>
            <CAA> <name> <ttl> <flags> <tag> <value>
            <MX> <name> <ttl> <priority> <mailserver>
        update <domain> <hash_id>
            <A|AAAA|CNAME|TXT> <name|@> <ttl> <ip|alias|text>
            <NS> <name|@> <ttl> <nameserver>
            <SSHFP> <name> <ttl> <algorithm> <fp_type> <fingerprint>
            <SRV> <name> <ttl> <priority> <weight> <port> <target>
            <TLSA> <name> <ttl> <certificate_usage> <selector> <matching_type> <hash>
            <CAA> <name> <ttl> <flags> <tag> <caavalue>
            <MX> <name> <ttl> <priority> <mailserver>

    domains
        list [-fn <name regex filter>]
        auth <domain> <language>
        detail <domain>
        update <domain> <admin_contact>
        transfer <domain> <auth>

Comments:
    filters are applied only to inline format
    parameters precedence is config_file > command_line > environment > defaults

`)

}

func main() {

    a24api := make(map[string]string)
    a24api_args := make(map[string]string)

// ================================================================================================================================================================
// PARSE ENVIRONMENT
// ================================================================================================================================================================

    a24api["endpoint"] = os.Getenv("A24API_ENDPOINT")
    a24api["token"] = os.Getenv("A24API_TOKEN")
    a24api["config"] = os.Getenv("A24API_CONFIG")

// ================================================================================================================================================================
// PARSE COMMAND-LINE
// ================================================================================================================================================================

    var params = os.Args[1:]
    var indexMax = len(params) - 1
    var indexUsedFlag = -1
    var posArgIndex = 0

    for index, element := range params {

        if index == indexUsedFlag {
            indexUsedFlag = -1
            continue
        } else {
            // print help and exit
            if (element == "-h" || element == "--help") {
                printHelp()
                os.Exit(0)
            // set config file
            } else if (element == "-c" || element == "--config") && (index < indexMax) && (a24api["service"] == "") {
                a24api["config"] = params[index + 1]
                indexUsedFlag = index + 1
            // set api endpoint
            } else if (element == "-e" || element == "--endpoint") && (index < indexMax) && (a24api["service"] == "") {
                a24api["endpoint"] = params[index + 1]
                indexUsedFlag = index + 1
            // set api token
            } else if (element == "-t" || element == "--token") && (index < indexMax) && (a24api["service"] == "") {
                a24api["token"] = params[index + 1]
                indexUsedFlag = index + 1
            // set output format
            } else if (element == "-f" || element == "--format") && (index < indexMax) && (a24api["service"] == "") {
                a24api["format"] = params[index + 1]
                indexUsedFlag = index + 1
            // set dial ip version to 4
            } else if (element == "-4") && (a24api["service"] == "") {
                a24api["dial"] = "tcp4"
            // set dial ip version to 6
            } else if (element == "-6") && (a24api["service"] == "") {
                a24api["dial"] = "tcp6"
            // set api service
            } else if (element == "dns" || element == "domain") && (a24api["service"] == "") {
                a24api["service"] = element
            // set api function
            } else if (element == "list" || element == "delete" || element == "create" || element == "update") && (a24api["service"] != "") {
                a24api["function"] = element
            // set name filter
            } else if (element == "-fn") && (index < indexMax) && (a24api["service"] != "") && (a24api["function"] != "") {
                a24api["filter-name"] = params[index + 1]
                indexUsedFlag = index + 1
            // set type filter
            } else if (element == "-ft") && (index < indexMax) && (a24api["service"] != "") && (a24api["function"] != "") {
                a24api["filter-type"] = params[index + 1]
                indexUsedFlag = index + 1
            // set value filter
            } else if (element == "-fv") && (index < indexMax) && (a24api["service"] != "") && (a24api["function"] != "") {
                a24api["filter-value"] = params[index + 1]
                indexUsedFlag = index + 1
            // set positional arguments
            } else if (a24api["service"] != "") && (a24api["function"] != "") {
                a24api_args["argument" + strconv.Itoa(posArgIndex)] = element
                posArgIndex++
            // exit on unexpected argument
            } else {
                fmt.Println("Unknown argument or argument out of order.")
                printHelp()
                os.Exit(1)
            }
        }
    }

// ================================================================================================================================================================
// LOAD DEFAULTS
// ================================================================================================================================================================

    if a24api["config"] == "" {
        var confPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
        if err != nil {
            fmt.Println(err)
        }
        a24api["config"] = confPath + "/" + Con_a24api_config
    }
    if a24api["endpoint"] == "" {
        a24api["endpoint"] = Con_a24api_endpoint
    }
    if a24api["token"] == "" {
        a24api["token"] = Con_a24api_token
    }
    if a24api["format"] == "" {
        a24api["format"] = Con_a24api_format
    }
    if a24api["filter-name"] == "" {
        a24api["filter-name"] = Con_a24api_name_regexp
    }
    if a24api["filter-type"] == "" {
        a24api["filter-type"] = Con_a24api_type_regexp
    }
    if a24api["filter-value"] == "" {
        a24api["filter-value"] = Con_a24api_value_regexp
    }
    if a24api["dial"] == "" {
        a24api["dial"] = Con_a24api_dial
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

// ================================================================================================================================================================
// PROCESS RESPONSE
// ================================================================================================================================================================

    if a24api["format"] == "json" {
        var pretty_json bytes.Buffer
        json.Indent(&pretty_json, a24api_response_body, "", "    ")
        fmt.Printf("%s\n", string(pretty_json.Bytes()))
    } else {
        // prepare regexp
        a24api_filter_name, _ := regexp.Compile(a24api["filter-name"])
        a24api_filter_type, _ := regexp.Compile(a24api["filter-type"])
        a24api_filter_value, _ := regexp.Compile(a24api["filter-value"])

        switch a24api["service"] {
            case "dns":
                if (a24api_response.StatusCode != 200) && (a24api_response.StatusCode != 204) {
                    fmt.Printf("%d %s\n", a24api_response.StatusCode, getCodeText(a24api_response.StatusCode, a24api["service"], a24api["function"]))
                    os.Exit(2)
                }
                switch a24api["function"] {
                    case "list":
                        // expected structure [ "domainA", "domainB" ]
                        if len(a24api_args) == 0 {
                            var structured_data []string
                            json.Unmarshal([]byte(a24api_response_body), &structured_data)
                            w := new(tabwriter.Writer)
                            w.Init(os.Stdout, 0, 8, 1, ' ', 0)
                            for _, element := range structured_data {
                                if a24api_filter_name.MatchString(element) {
                                    fmt.Fprintf(w, "%s\n", element)
                                }
                            }
                            w.Flush()
                        } else {
                        // expected structure [ { "variableA": "value", "variableB": "value" }, { "variableA": "value", "variableB": "value" } ]
                            var structured_data []map[string]interface{}
                            json.Unmarshal([]byte(a24api_response_body), &structured_data)
                            //fmt.Printf("%v\n", structured_data)
                            w := new(tabwriter.Writer)
                            w.Init(os.Stdout, 0, 8, 1, ' ', 0)
                            for _, element := range structured_data {
                                if a24api_filter_type.MatchString(element["type"].(string)) {
                                    switch element["type"].(string) {
                                        case "A", "AAAA":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["ip"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%s\n", a24api["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["ip"].(string))
                                            }
                                        case "CNAME":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["alias"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%s\n", a24api["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["alias"].(string))
                                            }
                                        case "TXT":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["text"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t\"%s\"\n", a24api["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["text"].(string))
                                            }
                                        case "NS":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["nameServer"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%s\n", a24api["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["nameServer"].(string))
                                            }
                                        case "SSHFP":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["text"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%g\t%g\t%s\n", a24api["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["algorithm"].(float64), element["fingerprintType"].(float64), element["text"].(string))
                                            }
                                        case "SRV":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["target"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%g\t%g\t%g\t%s\n", a24api["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["priority"].(float64), element["weight"].(float64), element["port"].(float64), element["target"].(string))
                                            }
                                        case "TLSA":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["hash"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%g\t%g\t%g\t%s\n", a24api["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["certificateUsage"].(float64), element["selector"].(float64), element["matchingType"].(float64), element["hash"].(string))
                                            }
                                        case "CAA":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["caaValue"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%g\t%s\t%s\n", a24api["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["flags"].(float64), element["tag"].(string), element["caaValue"].(string))
                                            }
                                        case "MX":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["mailserver"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%g\t%s\n", a24api["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["priority"].(float64), element["mailserver"].(string))
                                            }
                                    }
                                }
                            }
                            w.Flush()
                        }
                    case "create", "update", "delete":
                        fmt.Printf("%d %s\n", a24api_response.StatusCode, getCodeText(a24api_response.StatusCode, a24api["service"], a24api["function"]))
                        os.Exit(0)
                }
        }
    }
}
