package main

import (
	"fmt"
	"log"
	"os"
	"time"

	getip "github.com/alicecarra/dynamic-dns-updater/internal"
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

	cloudflareconfigs := cloudflare.CloudflareBasicConfigs{
		ZoneIdentifier: zone_identifier,
		AuthKey:        auth_key,
		AuthEmail:      auth_email,
	}

	dnsrecords := dns.Getdns(cloudflareconfigs)
	fmt.Println("Avalible Records")
	fmt.Printf("%+v\n", dnsrecords)


	for {
		dnsrecords := dns.Getdns(cloudflareconfigs)

		actual_ip, err := getip.Getip()
		if err != nil {
			log.Fatal(err)
		}

		for _, record := range dnsrecords {

			// check for record name specified on env
			if record.Name != recordname{
				continue 
			}

			if record.Content != actual_ip {
				log.Printf("Updating IP to %s\n", actual_ip)

				update_time := time.Now().Format(time.RFC822Z)

				dnsconfig := dns.DNSUpdateConfig{
					Name:     recordname,
					IP:       actual_ip,
					Comment:  fmt.Sprintf("last updated in %s", update_time),
					Proxied:  record.Proxied,
					Dns_type: "A",
					Zone_id:  record.Zone_id,
					Id:       record.Id,
				}

				fmt.Printf("%+v", dnsconfig)
		
				err = dns.Updatedns(dnsconfig, cloudflareconfigs)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		log.Println("IP not changed")
		time.Sleep(300 * time.Second)
	}
}

