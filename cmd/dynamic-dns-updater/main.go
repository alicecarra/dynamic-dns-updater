package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IP struct {
    Query string
}

func main() {

    m := getip()
    fmt.Print(m)

}

func getip() string {
    req, err := http.Get("http://ip-api.com/json/")
    if err != nil {
        return err.Error()
    }
    defer req.Body.Close()

    body, err := io.ReadAll(req.Body)
    if err != nil {
        return err.Error()
    }
    var ip IP
    json.Unmarshal(body, &ip)
    // fmt.Print(ip.Query)
    return ip.Query
}