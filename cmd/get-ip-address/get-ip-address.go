package main

import (
	"log"
	"os"

	"github.com/alicecarra/dynamic-dns-updater/pkg/cloudflare"
	"github.com/alicecarra/dynamic-dns-updater/pkg/cloudflare/dns"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("env file not loaded")
	}

	recordname := os.Getenv("RECORD_NAME")
	if recordname == "" {
		log.Fatal("Missing RECORD_NAME in env")
	}

	zone_identifier := os.Getenv("ZONE_IDENTIFIER")
	if zone_identifier == "" {
		log.Fatal("Missing ZONE_IDENTIFIER in env")
	}

	auth_key := os.Getenv("AUTH_KEY")
	if auth_key == "" {
		log.Fatal("Missing AUTH_KEY in env")
	}

	

	cloudflareconfigs := cloudflare.CloudflareBasicConfigs{
		ZoneIdentifier:    zone_identifier,
		AuthtorizationKey: auth_key,
	}

	dnsrecords, err := dns.Getdns(cloudflareconfigs)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Avalible Records")
	// fmt.Printf("%+v\n", dnsrecords)

	for _, record := range dnsrecords {

		// check for record name specified on env
		if record.Name != recordname{
			continue 
		}

			log.Printf("%s - %s - %s\n", record.Name, record.Content, record.Comment)
	}

}