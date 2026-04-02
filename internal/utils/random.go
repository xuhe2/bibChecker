package utils

import (
	"math/rand"
)

var userAgents = [...]string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
}

var acceptHeaders = [...]string{
	"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
}

var acceptLanguages = [...]string{
	"en-US,en;q=0.9",
	"zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
}

var secChUaList = [...]string{
	`"Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"`,
	`"Chromium";v="144", "Not-A.Brand";v="24", "Google Chrome";v="144"`,
}

var secChUaFullVersionList = [...]string{
	`"Chromium";v="146.0.7680.165", "Not-A.Brand";v="24.0.0.0", "Google Chrome";v="146.0.7680.165"`,
	`"Chromium";v="144.0.3425.32", "Not-A.Brand";v="24.0.0.0", "Google Chrome";v="144.0.3425.32"`,
}

var secChUaPlatform = [...]string{
	`"macOS"`,
	`"Windows"`,
}

func init() {
	rand.Shuffle(len(userAgents), func(i, j int) {
		userAgents[i], userAgents[j] = userAgents[j], userAgents[i]
	})
}

func UserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func Accept() string {
	return acceptHeaders[rand.Intn(len(acceptHeaders))]
}

func AcceptLanguage() string {
	return acceptLanguages[rand.Intn(len(acceptLanguages))]
}

func SecChUa() string {
	return secChUaList[rand.Intn(len(secChUaList))]
}

func SecChUaFullVersionList() string {
	return secChUaFullVersionList[rand.Intn(len(secChUaFullVersionList))]
}

func SecChUaPlatform() string {
	return secChUaPlatform[rand.Intn(len(secChUaPlatform))]
}
