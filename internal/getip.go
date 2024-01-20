package getip

import (
	"encoding/json"
	"io"
	"net/http"
)

type GetIpResponse struct {
	Query string
}

func Getip() (string, error) {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	var ip GetIpResponse

	err = json.Unmarshal(body, &ip)
	if err != nil {
		return "", err
	}
	return ip.Query, nil
}