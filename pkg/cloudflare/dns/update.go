package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/alicecarra/dynamic-dns-updater/pkg/cloudflare"
)

type DNSUpdateConfig struct {
	Name    	string `json:"name"`
	IP      	string `json:"content"`
	Comment 	string `json:"comment"`
	Proxied		bool   `json:"proxied"`
	Dns_type 	string `json:"type"`
	Zone_id		string `json:"zoneId"`
	Id			string `json:"id"`
}

func Updatedns(dnsupdateconfig DNSUpdateConfig, cloudflareconfigs cloudflare.CloudflareBasicConfigs) error {
	data, err := json.Marshal(dnsupdateconfig)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	apiurl := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", dnsupdateconfig.Zone_id, dnsupdateconfig.Id)
	request, err := http.NewRequest(http.MethodPut, apiurl, bytes.NewBuffer(data))
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
