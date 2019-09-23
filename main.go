package main

import (
    "os"
    "bytes"
    "encoding/json"
//    "fmt"
    "io/ioutil"
    "net/http"
    "path/filepath"
    "log"
    "strconv"
)

const (
    Con_a24api_endpoint = "https://sandboxapi.active24.com"
    Con_a24api_token = "123456qwerty-ok"
    Con_a24api_config = "a24api-conf.json"
)

func printHelp() {
    log.Println(`Usage: a24api [options] <service> <function> [parameters]

Options:
    -c|--config <path>            Path to config file. Default is a24api-conf.json.
    -e|--endpoint <url>           Active24 REST API url.
    -t|--token <token>            Active24 REST API token.

Services, functions and parameters:
    dns
        list [-fn <name regex filter>]
        list <domain> [-ft <type regex filter>] [-fn <name regex filter>] [-fv <value regex filter>]
        delete <domain> <record id>
        create <domain> <...>
            <A|AAAA|CNAME|TXT> <name|@> <ttl> <value>

        update <domain> <record id> <...>
            <A|AAAA|CNAME|TXT> <name|@> <ttl> <value>

    domain
        list
`)

}

func main() {

    a24api := make(map[string]string)

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
            log.Printf("Command-line argument already used: %d (of %d) = %s", index, indexMax, element)
            continue
        } else {
            log.Printf("Command-line argument to process: %d (of %d) = %s", index, indexMax, element)
            // print help and exit
            if (element == "-h" || element == "--help") {
                printHelp()
                os.Exit(0)
            // set config file
            } else if (element == "-c" || element == "--config") && (index < indexMax) && (a24api["service"] == "") {
                a24api["config"] = params[index + 1]
                indexUsedFlag = index + 1
                log.Printf("config set to: %s", params[index + 1])
            // set api endpoint
            } else if (element == "-e" || element == "--endpoint") && (index < indexMax) && (a24api["service"] == "") {
                a24api["endpoint"] = params[index + 1]
                indexUsedFlag = index + 1
                log.Printf("endpoint set to: %s", params[index + 1])
            // set api token
            } else if (element == "-t" || element == "--token") && (index < indexMax) && (a24api["service"] == "") {
                a24api["config"] = params[index + 1]
                indexUsedFlag = index + 1
                log.Printf("token set to: %s", params[index + 1])
            // set api service
            } else if (element == "dns" || element == "domain") && (a24api["service"] == "") {
                a24api["service"] = element
                log.Printf("service set to: %s", params[index])
            // set api function
            } else if (element == "list" || element == "delete" || element == "create" || element == "update") && (a24api["service"] != "") {
                a24api["function"] = element
                log.Printf("function set to: %s", params[index])
            // set name filter
            } else if (element == "-fn") && (index < indexMax) && (a24api["service"] != "") && (a24api["function"] != "") {
                a24api["filter-name"] = params[index + 1]
                log.Printf("filter-name set to: %s", params[index + 1])
            // set type filter
            } else if (element == "-ft") && (index < indexMax) && (a24api["service"] != "") && (a24api["function"] != "") {
                a24api["filter-type"] = params[index + 1]
                log.Printf("filter-type set to: %s", params[index + 1])
            // set value filter
            } else if (element == "-fv") && (index < indexMax) && (a24api["service"] != "") && (a24api["function"] != "") {
                a24api["filter-value"] = params[index + 1]
                log.Printf("filter-value set to: %s", params[index + 1])
            // set positional arguments
            } else if (a24api["service"] != "") && (a24api["function"] != "") {

                switch a24api["service"] {
                    case "dns":
                        // first argument should be always domain
                        if posArgIndex == 0 {
                            a24api["domain"] = element
                        }
                    default:
                        a24api["argument" + strconv.Itoa(posArgIndex)] = element
                        log.Printf("argument%s: %s", strconv.Itoa(posArgIndex), element)
                }
                posArgIndex++
            // exit on unexpected argument
            } else {
                log.Println("Unknown argument or argument out of order.")
                printHelp()
                log.Fatal()
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

// ================================================================================================================================================================
// LOAD CONFIG FILE
// ================================================================================================================================================================

// ================================================================================================================================================================
// CHECK INPUT DATA
// ================================================================================================================================================================

    if a24api["service"] == "" || a24api["function"] == "" {
        log.Println("Service of function not provided.\n")
        printHelp()
        log.Fatal()
    }

// ================================================================================================================================================================
// PROCESS
// ================================================================================================================================================================

//    a24api_request_body, err := json.Marshal(map[string]string{
//        "test": "test",
//        "test1": "test1",
//    })

    a24api_request_body := make(map[string]string)

    switch a24api["service"] {
        case "dns":
            switch a24api["function"] {
                case "list":
                    if a24api["domain"] == "" {
                        a24api["endpoint-uri"] = "/dns/domains/v1"
                        a24api["endpoint-method"] = "GET"
                    } else {
                        a24api["endpoint-uri"] = "/dns/" + a24api["domain"] + "/records/v1"
                        a24api["endpoint-method"] = "GET"
                    }
                case "create":
                    a24api["endpoint-uri"] = "/dns/" + a24api["domain"] + "/" + a24api["record-type"] + "/v1"
                    a24api["endpoint-method"] = "POST"
                    a24api_request_body[""] = ""
                case "update":
                    a24api["endpoint-uri"] = "/dns/" + a24api["domain"] + "/" + a24api["record-type"] + "/v1"
                    a24api["endpoint-method"] = "PUT"
                case "delete":
                    a24api["endpoint-uri"] = "/dns/" + a24api["domain"] + "/" + a24api["record-id"] + "/v1"
                    a24api["endpoint-method"] = "DELETE"
                default:
                log.Fatalln("Unknown function.")
            }
        default:
            log.Fatalln("Unknown service.")
    }

    a24api_client := &http.Client{}

    a24api_request_body_json, err := json.Marshal(a24api_request_body)
    if err != nil {
        log.Fatalln(err)
    }

    a24api_request, err := http.NewRequest(a24api["endpoint-method"], a24api["endpoint"] + a24api["endpoint-uri"], bytes.NewBuffer(a24api_request_body_json))
//    a24api_request, err := http.NewRequest(a24api["endpoint-method"], a24api["endpoint"] + a24api["endpoint-uri"], nil)
    if err != nil {
        log.Fatalln(err)
    }

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

    log.Println(string(a24api_response_body))

}
