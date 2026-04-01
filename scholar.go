package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {
	proxy := "http://127.0.0.1:7890"
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		fmt.Println("Invalid proxy URL:", err)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	baseURL := "https://scholar.google.com/scholar"
	params := url.Values{
		"hl":      {"zh-CN"},
		"as_sdt":  {"0,5"},
		"q":       {"Adaptive prototype learning for few-shot semantic segmentation"},
		"btnG":    {""},
	}

	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response length:", len(body))

	// 写入HTML文件方便调试
	filename := "scholar_output.html"
	if err := os.WriteFile(filename, body, 0644); err != nil {
		fmt.Println("Failed to write file:", err)
		return
	}
	fmt.Println("HTML saved to:", filename)
}