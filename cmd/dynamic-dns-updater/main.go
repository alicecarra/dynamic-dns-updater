package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
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

    auth_email := os.Getenv("AUTH_EMAIL")
    if auth_email == "" {
        log.Fatal("Missing AUTH_EMAIL in env")
    }

    auth_key := os.Getenv("AUTH_KEY")
    if auth_key == "" {
        log.Fatal("Missing AUTH_KEY in env")
    }

    cloudflareconfigs := CloudflareConfigs{
    	Identifier:     identifier,
    	ZoneIdentifier: zone_identifier,
    	AuthKey:        auth_key,
    	AuthEmail:      auth_email,
    }

    OUTER:
    for {
        dns_ip, err := net.LookupIP(recordname)
        if err != nil {
            log.Fatal(err)
        }

        actual_ip, err := getip()
        if err != nil {
            log.Fatal(err)
        }

        for _, ip := range dns_ip {
    
            if ip.String() == actual_ip { 
                log.Println("IP not changed")

                time.Sleep(300 * time.Second)

                continue OUTER
            }
    
        }

        log.Printf("Updating IP to %s\n", actual_ip)



        update_time := time.Now().Format(time.RFC822Z)
        
        dnsconfig := DNSConfig{
            Name:    recordname,
            IP:      actual_ip,
            Comment: fmt.Sprintf("last updated in %s", update_time),
            Proxied: false,
            DnsType: "A",
        }


        err = updatedns(dnsconfig, cloudflareconfigs)
    
        if err != nil {
            log.Fatal(err)
        }
    }
}


type CloudflareConfigs struct {
    Identifier string
    ZoneIdentifier string
    AuthKey string
    AuthEmail string
}

type DNSConfig struct {
    Name string `json:"name"`
    IP string `json:"content"`
    Comment string `json:"comment"`
    Proxied bool  `json:"proxied"`
    DnsType string `json:"type"`

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

func updatedns(dnsconfig DNSConfig, cloudflareconfigs CloudflareConfigs) error {
    data, err := json.Marshal(dnsconfig)
    if err != nil {
        return err
    }

    fmt.Println(string(data))

    apiurl := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", cloudflareconfigs.ZoneIdentifier, cloudflareconfigs.Identifier)
    request,    err := http.NewRequest(http.MethodPut, apiurl, bytes.NewBuffer(data))
    if err != nil {
        return err
    }

    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("X-Auth-Email", cloudflareconfigs.AuthEmail)
    request.Header.Set("X-Auth-Key", cloudflareconfigs.AuthKey)


    
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return err
    }

    responseData, err := io.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(responseData))

    defer response.Body.Close()


    return nil

}