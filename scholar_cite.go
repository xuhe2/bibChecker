package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {
	// 读取HTML文件
	html, err := os.ReadFile("scholar_output.html")
	if err != nil {
		fmt.Println("Failed to read file:", err)
		return
	}

	// 提取论文ID - 匹配 data-cid 属性
	re := regexp.MustCompile(`data-cid="([^"]+)"`)
	matches := re.FindAllSubmatch(html, -1)

	fmt.Println("找到", len(matches), "篇论文的引用链接")
	fmt.Println("========================================")

	// 使用代理
	proxy := "http://127.0.0.1:7890"
	proxyURL, _ := url.Parse(proxy)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// 提取唯一的论文ID
	uniqueIDs := make(map[string]bool)
	for _, m := range matches {
		id := string(m[1])
		if !uniqueIDs[id] {
			uniqueIDs[id] = true
		}
	}

	fmt.Printf("共有 %d 篇唯一论文\n\n", len(uniqueIDs))

	// 只处理前3篇进行测试
	count := 0
	for id := range uniqueIDs {
		if count >= 3 {
			break
		}
		count++

		fmt.Printf("[%d] 论文ID: %s\n", count, id)

		// Step 1: 请求引用页面，获取 scisdr 和 scisig 参数
		citeURL := fmt.Sprintf("https://scholar.google.com/scholar?q=info:%s:scholar.google.com/&output=cite&scirp=0&hl=zh-CN", id)

		fmt.Printf("    Step 1 - Cite URL: %s\n", citeURL)

		req, _ := http.NewRequest("GET", citeURL, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("    Error: %v\n\n", err)
			continue
		}

		citeBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		citeContent := string(citeBody)

		// 保存引用页面用于调试
		citeFilename := fmt.Sprintf("cite_%s.html", id)
		os.WriteFile(citeFilename, []byte(citeContent), 0644)
		fmt.Printf("    Cite page saved to: %s\n", citeFilename)

		// 从引用页面提取 scisdr 和 scisig
		reScisdr := regexp.MustCompile(`scisdr=([^&]+)`)
		reScisig := regexp.MustCompile(`scisig=([^&]+)`)

		scisdrMatch := reScisdr.FindStringSubmatch(citeContent)
		scisigMatch := reScisig.FindStringSubmatch(citeContent)

		if scisdrMatch == nil || scisigMatch == nil {
			fmt.Printf("    Failed to extract scisdr/scisig\n")
			// 打印部分内容用于调试
			if len(citeContent) > 300 {
				fmt.Printf("    Content preview: %s\n", citeContent[:300])
			} else {
				fmt.Printf("    Content: %s\n", citeContent)
			}
			fmt.Println("----------------------------------------")
			continue
		}

		scisdr := scisdrMatch[1]
		scisig := scisigMatch[1]
		fmt.Printf("    scisdr: %s\n", scisdr)
		fmt.Printf("    scisig: %s\n", scisig)

		// Step 2: 使用提取的参数请求 BibTeX
		bibURL := fmt.Sprintf("https://scholar.googleusercontent.com/scholar.bib?q=info:%s:scholar.google.com/&output=citation&scisdr=%s&scisig=%s&scisf=4&ct=citation&cd=-1&hl=zh-CN",
			id, scisdr, scisig)

		fmt.Printf("    Step 2 - BibTeX URL: %s\n", bibURL)

		req2, _ := http.NewRequest("GET", bibURL, nil)
		req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

		resp2, err := client.Do(req2)
		if err != nil {
			fmt.Printf("    Error: %v\n\n", err)
			continue
		}

		bibBody, _ := io.ReadAll(resp2.Body)
		resp2.Body.Close()
		bibContent := string(bibBody)

		// 保存 BibTeX 文件
		bibFilename := fmt.Sprintf("%s.bib", id)
		if err := os.WriteFile(bibFilename, []byte(bibContent), 0644); err != nil {
			fmt.Printf("    Failed to write BibTeX: %v\n", err)
		} else {
			fmt.Printf("    BibTeX saved to: %s\n", bibFilename)
		}

		// 打印 BibTeX 内容
		if strings.Contains(bibContent, "@") {
			fmt.Printf("    BibTeX content:\n%s\n", bibContent)
		} else {
			fmt.Printf("    BibTeX response: %s\n", bibContent)
		}
		fmt.Println("----------------------------------------")
	}
}