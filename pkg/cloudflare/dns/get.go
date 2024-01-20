package dns

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alicecarra/dynamic-dns-updater/pkg/cloudflare"
)

/// Used only for parsing response from API
type DnsRecordsResponse struct {
	Result []DNSRecord
}


func Getdns(cloudflareconfigs cloudflare.CloudflareBasicConfigs) []DNSRecord {
	apiurl := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", cloudflareconfigs.ZoneIdentifier)

	request, err := http.NewRequest(
		http.MethodGet,
		apiurl,
		nil,
	)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Email", cloudflareconfigs.AuthEmail)
	request.Header.Set("X-Auth-Key", cloudflareconfigs.AuthKey)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatalf("error sending HTTP request: %v", err)
	}

	var dnsrecords DnsRecordsResponse
	d := json.NewDecoder(response.Body)
	if err := d.Decode(&dnsrecords); err != nil {
		log.Fatalf("error deserializing dns data: %v", err)
	}

	return dnsrecords.Result
}



