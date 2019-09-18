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

Services nad functions:
    dns
        list
        create
        update
        delete
    domain
        list
        create
        update
        delete
`)

}

func main() {
// ================================================================================================================================================================
// PARSE COMMAND-LINE
// ================================================================================================================================================================

    var cmd_a24api_config = ""
    var cmd_a24api_service = ""
    var cmd_a24api_function = ""

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




// ================================================================================================================================================================
// PARSE CONFIG FILE
// ================================================================================================================================================================


// ================================================================================================================================================================
// PROCESS
// ================================================================================================================================================================

    a24api_headers = {
            "Accept": "application/json",
            "Authorization": "Bearer " + apikey
        }


    httpClient := &http.Client{}
    req, err := http.NewRequest("GET", run_a24api_endpoint, nil)


    response, err := http.Get("https://httpbin.org/ip")
    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
    }
    jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
    jsonValue, _ := json.Marshal(jsonData)
    response, err = http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonValue))
    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
    }
}
