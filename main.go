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
    Con_a24api_key = "123456qwerty-ok"
    Con_a24api_config = "a24api-conf.json"
)

func printHelp() {
    log.Println(`Usage: a24api [options] <service> <function> [parameters]

Options:
    -c|--config <path>            Path to config file. Default is a24api-conf.json.
    -e|--endpoint <url>           Active24 REST API url.
    -k|--key <key>                Active24 REST API key.

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
// ================================================================================================================================================================
// PARSE COMMAND-LINE
// ================================================================================================================================================================

    var cmd_a24api_endpoint = ""
    var cmd_a24api_key = ""

    var cmd_a24api_config = ""
//    var cmd_a24api_service = ""
//    var cmd_a24api_function = ""

    var argCnt = len(os.Args[1:])
    for index, element := range os.Args[1:] {
        log.Printf("Command-line arguments: %d (of %d) = %s", index, argCnt, element)
    }

// ================================================================================================================================================================
// PARSE ENVIRONMENT
// ================================================================================================================================================================

    var env_a24api_endpoint = os.Getenv("A24API_ENDPOINT")
    var env_a24api_key = os.Getenv("A24API_KEY")
    var env_a24api_config = os.Getenv("A24API_CONFIG")

// ================================================================================================================================================================
// HANDLE "env" -> "cmd" -> "Con" PRECEDENCE
// ================================================================================================================================================================

    // run_a24api_config
    if env_a24api_config != "" {
        var run_a24api_config = env_a24api_config
        if _, err := os.Stat(run_a24api_config); err != nil {
            log.Fatal("Environment variable A24API_CONFIG is set, but provided file does not exists.")
        }
    } else if cmd_a24api_config != "" {
        var run_a24api_config = cmd_a24api_config
        if _, err := os.Stat(run_a24api_config); err != nil {
            log.Fatal("Command-line argument --config is set, but provided file does not exists.")
        }
    } else {
        var run_a24api_config, _ = filepath.Abs(filepath.Dir(os.Args[0]))
        run_a24api_config = run_a24api_config + "/" + Con_a24api_config
        if _, err := os.Stat(run_a24api_config); err != nil {
            log.Println("Config file does not exists.")
        }
    }
    // run_a24api_endpoint
    var run_a24api_endpoint = ""
    if env_a24api_endpoint != "" {
        run_a24api_endpoint = env_a24api_endpoint
    } else if cmd_a24api_endpoint != "" {
        run_a24api_endpoint = cmd_a24api_endpoint
    } else {
        run_a24api_endpoint = Con_a24api_endpoint
    }
    // run_a24api_key
    var run_a24api_key = ""
    if env_a24api_key != "" {
        run_a24api_key = env_a24api_key
    } else if cmd_a24api_key != "" {
        run_a24api_key = cmd_a24api_key
    } else {
        run_a24api_key = Con_a24api_key
    }

// ================================================================================================================================================================
// PARSE CONFIG FILE
// ================================================================================================================================================================


// ================================================================================================================================================================
// PROCESS
// ================================================================================================================================================================

    a24api_client := &http.Client{}
    a24api_request, err := http.NewRequest("GET", run_a24api_endpoint, nil)
    if err != nil {
        log.Fatalln(err)
    }

    a24api_request.Header.Set("Accept", "application/json")
    a24api_request.Header.Set("Authorization", "Bearer " + run_a24api_key)

    a24api_response, err := a24api_client.Do(a24api_request)
    if err != nil {
        log.Fatalln(err)
    }

    defer a24api_response.Body.Close()

    a24api_body, err := ioutil.ReadAll(a24api_response.Body)
    if err != nil {
        log.Fatalln(err)
    }

    log.Println(string(a24api_body))



//    response, err := http.Get("https://httpbin.org/ip")
//    if err != nil {
//        fmt.Printf("The HTTP request failed with error %s\n", err)
//    } else {
//        data, _ := ioutil.ReadAll(response.Body)
//        fmt.Println(string(data))
//    }
//    jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
//    jsonValue, _ := json.Marshal(jsonData)
//    response, err = http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonValue))
//    if err != nil {
//        fmt.Printf("The HTTP request failed with error %s\n", err)
//    } else {
//        data, _ := ioutil.ReadAll(response.Body)
//        fmt.Println(string(data))
//    }
}
