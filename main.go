package main

import (
    "os"
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "path/filepath"
    "log"
    "strconv"
    "strings"
//    "text/tabwriter"
)

const (
    Con_a24api_endpoint = "https://sandboxapi.active24.com"
    Con_a24api_token = "123456qwerty-ok"
    Con_a24api_config = "a24api-conf.json"
    Con_a24api_format = "inline"
)

type a24api_config_file_t struct {
    a24api_endpoint string
    a24api_token string
}

func printHelp() {
    log.Println(`Usage: a24api [options] <service> <function> [parameters]

Options:
    -c|--config <path>            Path to config file. Default is a24api-conf.json.
    -e|--endpoint <url>           Active24 REST API url.
    -t|--token <token>            Active24 REST API token.
    -f|--format <json|inline>     Output format (default: inline).

Services, functions and parameters:
    dns
        list [-fn <name regex filter>]
        list <domain> [-ft <type regex filter>] [-fn <name regex filter>] [-fv <value regex filter>]
        delete <domain> <hash_id>
        create <domain>
            <A|AAAA|CNAME|TXT> <name|@> <ttl> <ip|alias|text>
            <SSHFP> <name> <ttl> <algorithm> <fp_type> <fingerprint>
            <SRV> <name> <ttl> <priority> <weight> <port> <target>
            <TLSA> <name> <ttl> <certificate_usage> <selector> <matching_type> <hash>
            <CAA> <name> <ttl> <flags> <tag> <value>
            <MX> <name> <ttl> <priority> <value>
        update <domain> <hash_id>
            <A|AAAA|CNAME|TXT> <name|@> <ttl> <ip|alias|text>
            <SSHFP> <name> <ttl> <algorithm> <fp_type> <fingerprint>
            <SRV> <name> <ttl> <priority> <weight> <port> <target>
            <TLSA> <name> <ttl> <certificate_usage> <selector> <matching_type> <hash>
            <CAA> <name> <ttl> <flags> <tag> <value>
            <MX> <name> <ttl> <priority> <value>
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
            // set api service
            } else if (element == "dns" || element == "domain") && (a24api["service"] == "") {
                a24api["service"] = element
            // set api function
            } else if (element == "list" || element == "delete" || element == "create" || element == "update") && (a24api["service"] != "") {
                a24api["function"] = element
            // set name filter
            } else if (element == "-fn") && (index < indexMax) && (a24api["service"] != "") && (a24api["function"] != "") {
                a24api["filter-name"] = params[index + 1]
            // set type filter
            } else if (element == "-ft") && (index < indexMax) && (a24api["service"] != "") && (a24api["function"] != "") {
                a24api["filter-type"] = params[index + 1]
            // set value filter
            } else if (element == "-fv") && (index < indexMax) && (a24api["service"] != "") && (a24api["function"] != "") {
                a24api["filter-value"] = params[index + 1]
            // set positional arguments
            } else if (a24api["service"] != "") && (a24api["function"] != "") {
                a24api_args["argument" + strconv.Itoa(posArgIndex)] = element
                posArgIndex++
            // exit on unexpected argument
            } else {
                log.Println("Unknown argument or argument out of order.")
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
            log.Fatalln(err)
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
        log.Println("Service or function not provided.")
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
                // expected arguments: 0=domain
                case "list":
                    if len(a24api_args) == 0 {
                        a24api["endpoint-uri"] = "/dns/domains/v1"
                        a24api["endpoint-method"] = "GET"
                    } else {
                        a24api["endpoint-uri"] = "/dns/" + a24api_args["argument0"] + "/records/v1"
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
                    a24api["record-type"] = strings.ToLower(a24api_args["argument" + strconv.Itoa(posArgOffset + 1)])
                    a24api["endpoint-uri"] = "/dns/" + a24api_args["argument0"] + "/" + a24api["record-type"] + "/v1"
                    switch a24api["record-type"] {
                        case "a", "aaaa":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=value
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["ip"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                        case "cname":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=value
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["alias"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                        case "txt":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=value
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["text"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                        case "sshfp":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=algorithm, 5(6)=fp_type, 6(7)=fingerprint
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["algorithm"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                            a24api_request_body["fingerprintType"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 5)]
                            a24api_request_body["text"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 6)]
                        case "srv":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=priority, 5(6)=weight, 6(7)=port, 7(8)=target
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["priority"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                            a24api_request_body["weight"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 5)]
                            a24api_request_body["port"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 6)]
                            a24api_request_body["target"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 7)]
                        case "tlsa":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=certificate_usage, 5(6)=selector, 6(7)=matching_type, 7(8)=hash
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["certificateUsage"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                            a24api_request_body["selector"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 5)]
                            a24api_request_body["matchingType"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 6)]
                            a24api_request_body["hash"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 7)]
                        case "caa":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=flags, 5(6)=tag, 6(7)=value
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["flags"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                            a24api_request_body["tag"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 5)]
                            a24api_request_body["caaValue"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 6)]
                        case "mx":
                            // expected arguments: 0=domain, 1(2)=type, 2(3)=name, 3(4)=ttl, 4(5)=priority, 5(6)=value
                            a24api_request_body["name"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 2)]
                            a24api_request_body["ttl"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 3)]
                            a24api_request_body["priority"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 4)]
                            a24api_request_body["mailserver"] = a24api_args["argument" + strconv.Itoa(posArgOffset + 5)]
                        default:
                            log.Fatalf("Unsupported dns type: %s.\n", a24api["record-type"])
                    }
                // expected arguments: 0=domain, 1=hash_id
                case "delete":
                    a24api["endpoint-uri"] = "/dns/" + a24api["domain"] + "/" + a24api["record-id"] + "/v1"
                    a24api["endpoint-method"] = "DELETE"
                    a24api_request_body["hashId"] = a24api_args["argument1"]
                default:
                    log.Fatalf("Unsupported function: %s.\n", a24api["function"])
            }
        default:
            log.Fatalf("Unsupported service: %s.", a24api["service"])
    }

// ================================================================================================================================================================
// MAKE REQUEST
// ================================================================================================================================================================

    a24api_client := &http.Client{}

    a24api_request_body_json, err := json.Marshal(a24api_request_body)
    if err != nil {
        log.Fatalln(err)
    }

    a24api_request, err := http.NewRequest(a24api["endpoint-method"], a24api["endpoint"] + a24api["endpoint-uri"], bytes.NewBuffer(a24api_request_body_json))
    if err != nil {
        log.Fatalln(err)
    }

    a24api_request.Header.Set("Content-type", "application/json")
    a24api_request.Header.Set("Accept", "application/json")
    a24api_request.Header.Set("Authorization", "Bearer " + a24api["token"])

    a24api_response, err := a24api_client.Do(a24api_request)
    if err != nil {
        log.Fatalln(err)
    }

    defer a24api_response.Body.Close()

    a24api_response_body, err := ioutil.ReadAll(a24api_response.Body)
    if err != nil {
        log.Fatalln(err)
    }

// ================================================================================================================================================================
// PROCESS RESPONSE
// ================================================================================================================================================================

//    if a24api_response.StatusCode == 200 {
//    var a24api_response_data a24api_response_data_t
//    json.Unmarshal([]byte(a24api_response_body), &a24api_response_data)

    if a24api["format"] == "json" {
        var out bytes.Buffer
        json.Indent(&out, a24api_response_body, "", "    ")
        fmt.Printf("%s\n", string(out.Bytes()))
    } else {
        var out bytes.Buffer
        json.Indent(&out, a24api_response_body, "", "    ")
        fmt.Printf("%s\n", string(out.Bytes()))
    }

}
