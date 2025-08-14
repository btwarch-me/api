package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func addCname (cname string,target string) error {
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	domain := os.Getenv("Domain")

	fullName := fmt.Sprintf("%s.%s",cname,domain)
	
	record := map[string]interface{}{
        "type":    "CNAME",
        "name":    fullName,
        "content": target,
        "ttl":     120,
        "proxied": false,
    }
	
	jsonData, _ := json.Marshal(record)
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	client := &http.Client{}

	resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
	
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        fmt.Println("CNAME record created:", string(body))
        return nil
    } else {
        return fmt.Errorf("failed to create CNAME: %s", string(body))
    }
}