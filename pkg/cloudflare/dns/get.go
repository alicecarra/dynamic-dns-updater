package dns

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alicecarra/dynamic-dns-updater/pkg/cloudflare"
)

/// Used only for parsing response from API
type DnsRecordsResponse struct {
	Result []DNSRecord
}


func Getdns(cloudflareconfigs cloudflare.CloudflareBasicConfigs) ([]DNSRecord, error) {
	apiurl := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", cloudflareconfigs.ZoneIdentifier)

	request, err := http.NewRequest(
		http.MethodGet,
		apiurl,
		nil,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	// request.Header.Set("X-Auth-Email", cloudflareconfigs.AuthEmail)
	// request.Header.Set("X-Auth-Key", cloudflareconfigs.AuthKey)
	request.Header.Set("Authorization", cloudflareconfigs.AuthtorizationKey)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	var dnsrecords DnsRecordsResponse
	d := json.NewDecoder(response.Body)
	if err := d.Decode(&dnsrecords); err != nil {
		return nil, err
	}

	return dnsrecords.Result, nil
}



