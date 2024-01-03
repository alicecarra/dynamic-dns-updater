package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)


func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    recordname := os.Getenv("RECORD_NAME")
    if recordname == "" {
        log.Fatal("Missing RECORD_NAME in env")
    }

    identifier := os.Getenv("IDENTIFIER")
    if identifier == "" {
        log.Fatal("Missing IDENTIFIER in env")
    }

    zone_identifier := os.Getenv("ZONE_IDENTIFIER")
    if zone_identifier == "" {
        log.Fatal("Missing ZONE_IDENTIFIER in env")
    }

    cloudflareconfigs := CloudflareConfigs{
    	Identifier:     identifier,
    	ZoneIdentifier: zone_identifier,
    }

    ip, err := getip()
    if err != nil {
        log.Fatal(err)
    }

    update_time := time.Now().Format(time.RFC822Z)
    
    dnsconfig := DNSConfig{
        Name:    recordname,
    	IP:      ip,
    	Comment: fmt.Sprintf("last updated in %s", update_time),
    	proxied: false,
    	DnsType: "A",
    }

    fmt.Printf("%+v\n%+v", cloudflareconfigs, dnsconfig)

    // updatedns(dnsconfig, cloudflareconfigs)

}


type CloudflareConfigs struct {
    Identifier string
    ZoneIdentifier string 
}

type DNSConfig struct {
    Name string
    IP string
    Comment string
    proxied bool
    DnsType string

}

type IP struct {
    Query string
}


func getip() (string, error) {
    req, err := http.Get("http://ip-api.com/json/")
    if err != nil {
        return "", err
    }
    defer req.Body.Close()

    body, err := io.ReadAll(req.Body)
    if err != nil {
        return "", err
    }
    var ip IP

    err = json.Unmarshal(body, &ip)
    if err != nil {
        return "", err
    }
    return ip.Query, nil
}
