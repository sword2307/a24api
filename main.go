package main

import (
    "os"
//    "bytes"
//    "encoding/json"
//    "fmt"
    "io/ioutil"
    "net/http"
    "path/filepath"
    "log"
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
        list <domain> [-t <type>] [-fn <name regex filter>] [-fv <value regex filter>]
        delete <domain> <record id>
        create
        update
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

    for index, element := range params {

        if index == indexUsedFlag {
            indexUsedFlag = -1
            log.Printf("Command-line argument already used: %d (of %d) = %s", index, indexMax, element)
            continue
        } else {
            log.Printf("Command-line argument to process: %d (of %d) = %s", index, indexMax, element)
            switch element {
                case "-c", "--config":
                    if (index < indexMax) && (a24api["service"] == "") {
                        a24api["config"] = params[index + 1]
                        indexUsedFlag = index + 1
                        log.Printf("config set to: %s", params[index + 1])
                    }
                case "-e", "--endpoint":
                    if (index < indexMax) && (a24api["service"] == "") {
                        a24api["endpoint"] = params[index + 1]
                        indexUsedFlag = index + 1
                        log.Printf("endpoint set to: %s", params[index + 1])
                    }
                case "-t", "--token":
                    if (index < indexMax) && (a24api["service"] == "") {
                        a24api["config"] = params[index + 1]
                        indexUsedFlag = index + 1
                        log.Printf("token set to: %s", params[index + 1])
                    }
                default:
                    log.Printf("Command-line arguments: %d (of %d) = %s", index, indexMax, element)
            }
        }
    }


// ================================================================================================================================================================
// LOAD DEFAULTS
// ================================================================================================================================================================

    if a24api["config"] == "" {
        var confPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
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
// PROCESS
// ================================================================================================================================================================

//    a24api_request_body, err := json.Marshal(map[string]string{
//        "test": "test",
//        "test1": "test1",
//    })

    a24api_client := &http.Client{}
//    a24api_request, err := http.NewRequest("GET", run_a24api_endpoint + "/dns/domains/v1", a24api_request_body)
    a24api_request, err := http.NewRequest("GET", a24api["endpoint"] + "/dns/domains/v1", nil)
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
