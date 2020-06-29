package main

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
    "strings"
    "a24api/lib"
)

const (
    C_A24ApiClient_Configfile string = "a24api-conf.json"
)

var (
    A24ApiClient                        *a24apiclient.T_A24ApiClient
    A24ApiClientConfig                  map[string]string
    A24ApiClientArgs                    map[string]string
    A24ApiClientFuncArgs                map[string]string

    A24ApiClientConfigArgs =            [...]string { "endpoint", "token", "network", "timeout" }
)

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
        list
        list <domain>
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
        list
        auth <domain> <language>
        detail <domain>
        update <domain> <admin_contact>
        transfer <domain> <auth>

Comments:
    parameters precedence is config_file > command_line > environment > defaults

`)

}

func main() {

    A24ApiClientConfig := make(map[string]string)
    A24ApiClientArgs := make(map[string]string)
    A24ApiClientFuncArgs := make(map[string]string)

// ================================================================================================================================================================
// PARSE ENVIRONMENT
// ================================================================================================================================================================

    A24ApiClientConfig["endpoint"] = os.Getenv("A24API_ENDPOINT")
    A24ApiClientConfig["token"] = os.Getenv("A24API_TOKEN")
    A24ApiClientConfig["config"] = os.Getenv("A24API_CONFIG")

// ================================================================================================================================================================
// PARSE COMMAND-LINE
// ================================================================================================================================================================

    var params = os.Args[1:]
    var indexMax = len(params) - 1
    var indexUsedFlag = -1
    var posFuncArgIndex = 0

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
            } else if (element == "-c" || element == "--config") && (index < indexMax) && (A24ApiClientConfig["service"] == "") {
                A24ApiClientConfig["config"] = params[index + 1]
                indexUsedFlag = index + 1
            // set api endpoint
            } else if (element == "-e" || element == "--endpoint") && (index < indexMax) && (A24ApiClientConfig["service"] == "") {
                A24ApiClientConfig["endpoint"] = params[index + 1]
                indexUsedFlag = index + 1
            // set api token
            } else if (element == "-t" || element == "--token") && (index < indexMax) && (A24ApiClientConfig["service"] == "") {
                A24ApiClientConfig["token"] = params[index + 1]
                indexUsedFlag = index + 1
            // set output format
            } else if (element == "-f" || element == "--format") && (index < indexMax) && (A24ApiClientConfig["service"] == "") {
                A24ApiClientArguments["format"] = params[index + 1]
                indexUsedFlag = index + 1
            // set network ip version to 4
            } else if (element == "-4") && (A24ApiClientConfig["service"] == "") {
                A24ApiClientConfig["network"] = "tcp4"
            // set network ip version to 6
            } else if (element == "-6") && (A24ApiClientConfig["service"] == "") {
                A24ApiClientConfig["network"] = "tcp6"
            // set api service
            } else if (element == "dns" || element == "domain") && (A24ApiClientConfig["service"] == "") {
                A24ApiClientArgs["service"] = element
            // set api function
            } else if (element == "list" || element == "delete" || element == "create" || element == "update") && (A24ApiClientConfig["service"] != "") {
                A24ApiClientArgs["function"] = element
            // set positional arguments
            } else if (A24ApiClientConfig["service"] != "") && (A24ApiClientConfig["function"] != "") {
                A24ApiClientFuncArgs[strconv.Itoa(posFuncArgIndex)] = element
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

    if A24ApiClientConfig["config"] == "" {
        var confPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
        if err != nil {
            fmt.Println(err)
        }
        A24ApiClientConfig["config"] = confPath + "/" + C_A24ApiClient_Configfile
    }

// ================================================================================================================================================================
// LOAD CONFIG FILE
// ================================================================================================================================================================

    configFile, err := os.Open(A24ApiClientConfig["config"])
    if err == nil {
        defer configFile.Close()
        configDataRaw, _ := ioutil.ReadAll(configFile)
        var configData map[string]string
        json.Unmarshal(configDataRaw, &configData)
        for _, lConfigArg := range A24ApiClientConfigArgs {
            if configData[lConfigArg] != "" {
                A24ApiClientConfig[lConfigArg] = configData[lConfigArg]
            }
        }
    }

// ================================================================================================================================================================
// CHECK INPUT DATA
// ================================================================================================================================================================

    if A24ApiClientConfig["service"] == "" || A24ApiClientConfig["function"] == "" {
        fmt.Println("Service or function not provided.")
        printHelp()
        os.Exit(1)
    }

// ================================================================================================================================================================
// INITIALIZE CLIENT
// ================================================================================================================================================================

    A24ApiClient := a24apiclient.NewA24ApiClient(A24ApiClientConfig)

// ================================================================================================================================================================
// PREPARE REQUEST
// ================================================================================================================================================================

    var posArgOffset = 0

    switch A24ApiClientArguments["service"] {
        case "dns":
            switch A24ApiClientArguments["function"] {
                case "list":
                    // expected arguments:
                    if A24ApiClientArguments["argument0"] == "" {
                        A24ApiResponseCode, A24ApiResponseBody, err := A24ApiClient.DnsListDomains()
                    // expected arguments: 0=domain
                    } else {
                        A24ApiResponseCode, A24ApiResponseBody, err := A24ApiClient.DnsListRecords(A24ApiClientArguments["argument0"])
                    }
                case "create", "update":
                    if A24ApiClient.Config["function"] == "create" {
                        A24ApiClient.Config["endpoint-method"] = "POST"
                    // expected arguments: 0=domain, 1=hash_id
                    } else {
                        A24ApiClient.Config["endpoint-method"] = "PUT"
                        a24api_request_body["hashId"] = a24api_args["argument1"]
                        posArgOffset = 1
                    }
                    // expected arguments: 0=domain, 1(2)=type
                    A24ApiClient.Config["domain"] = a24api_args["argument0"]
                    A24ApiClient.Config["record-type"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 1)]
                    A24ApiClient.Config["endpoint-uri"] = "/dns/" + A24ApiClient.Config["domain"] + "/" + strings.ToLower(A24ApiClient.Config["record-type"]) + "/v1"
                    switch A24ApiClient.Config["record-type"] {
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
                            fmt.Printf("Unsupported dns type: %s.\n", A24ApiClient.Config["record-type"])
                            os.Exit(1)
                    }
                // expected arguments: 0=domain, 1=hash_id
                case "delete":
                    A24ApiClient.Config["domain"] = a24api_args["argument0"]
                    A24ApiClient.Config["endpoint-uri"] = "/dns/" + A24ApiClient.Config["domain"] + "/" + a24api_args["argument1"] + "/v1"
                    A24ApiClient.Config["endpoint-method"] = "DELETE"
                    a24api_request_body["hashId"] = a24api_args["argument1"]
                default:
                    fmt.Printf("Unsupported function: %s.\n", A24ApiClient.Config["function"])
                    os.Exit(1)
            }
        default:
            fmt.Printf("Unsupported service: %s.", A24ApiClient.Config["service"])
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
            }).Dial(A24ApiClient.Config["dial"], addr)
        }),
    }

    a24api_client := &http.Client{Transport: a24api_transport}

    a24api_request_body_json, err := json.Marshal(a24api_request_body)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    a24api_request, err := http.NewRequest(A24ApiClient.Config["endpoint-method"], A24ApiClient.Config["endpoint"] + A24ApiClient.Config["endpoint-uri"], bytes.NewBuffer(a24api_request_body_json))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    a24api_request.Header.Set("Content-type", "application/json")
    a24api_request.Header.Set("Accept", "application/json")
    a24api_request.Header.Set("Authorization", "Bearer " + A24ApiClient.Config["token"])

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

    if A24ApiClientArguments["format"] == "json" {
        var pretty_json bytes.Buffer
        json.Indent(&pretty_json, a24api_response_body, "", "    ")
        fmt.Printf("%s\n", string(pretty_json.Bytes()))
    } else {

        switch A24ApiClient.Config["service"] {
            case "dns":
                if (a24api_response.StatusCode != 200) && (a24api_response.StatusCode != 204) {
                    fmt.Printf("%d %s\n", a24api_response.StatusCode, getCodeText(a24api_response.StatusCode, A24ApiClient.Config["service"], A24ApiClient.Config["function"]))
                    os.Exit(2)
                }
                switch A24ApiClient.Config["function"] {
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
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%s\n", A24ApiClient.Config["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["ip"].(string))
                                            }
                                        case "CNAME":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["alias"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%s\n", A24ApiClient.Config["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["alias"].(string))
                                            }
                                        case "TXT":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["text"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t\"%s\"\n", A24ApiClient.Config["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["text"].(string))
                                            }
                                        case "NS":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["nameServer"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%s\n", A24ApiClient.Config["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["nameServer"].(string))
                                            }
                                        case "SSHFP":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["text"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%g\t%g\t%s\n", A24ApiClient.Config["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["algorithm"].(float64), element["fingerprintType"].(float64), element["text"].(string))
                                            }
                                        case "SRV":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["target"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%g\t%g\t%g\t%s\n", A24ApiClient.Config["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["priority"].(float64), element["weight"].(float64), element["port"].(float64), element["target"].(string))
                                            }
                                        case "TLSA":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["hash"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%g\t%g\t%g\t%s\n", A24ApiClient.Config["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["certificateUsage"].(float64), element["selector"].(float64), element["matchingType"].(float64), element["hash"].(string))
                                            }
                                        case "CAA":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["caaValue"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%g\t%s\t%s\n", A24ApiClient.Config["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["flags"].(float64), element["tag"].(string), element["caaValue"].(string))
                                            }
                                        case "MX":
                                            if a24api_filter_name.MatchString(element["name"].(string)) && a24api_filter_value.MatchString(element["mailserver"].(string)) {
                                                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%g\t%s\n", A24ApiClient.Config["domain"], element["hashId"].(string), element["type"].(string), element["name"].(string), element["ttl"].(float64), element["priority"].(float64), element["mailserver"].(string))
                                            }
                                    }
                                }
                            }
                            w.Flush()
                        }
                    case "create", "update", "delete":
                        fmt.Printf("%d %s\n", a24api_response.StatusCode, getCodeText(a24api_response.StatusCode, A24ApiClient.Config["service"], A24ApiClient.Config["function"]))
                        os.Exit(0)
                }
        }
    }
}
