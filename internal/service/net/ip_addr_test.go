package net_test

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/AmbitiousJun/live-server/internal/util/https"
)

func TestReg(t *testing.T) {
	url := `https://ipchaxun.com/172.104.64.164/`
	asnReg := regexp.MustCompile(`<span class="name">归属地：</span><span class="value">(.*)<a href="[^"]*" target="_blank" rel="nofollow">(.*)</a>(.*)</span>`)
	asn1Reg := regexp.MustCompile(`<span class="name">归属地：</span><span class="value">(.*)</span>`)
	providerReg := regexp.MustCompile(`<label><span class="name">运营商：</span><span class="value">(.*)</span></label>`)
	header := make(http.Header)
	header.Set("User-Agent", "libmpv")
	resp, err := https.Get(url).Header(header).Do()
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	if asnReg.MatchString(body) {
		log.Println(asnReg.FindStringSubmatch(body)[1])
		log.Println(asnReg.FindStringSubmatch(body)[2])
		log.Println(asnReg.FindStringSubmatch(body)[3])
	}
	if asn1Reg.MatchString(body) {
		log.Println(asn1Reg.FindStringSubmatch(body)[1])
	}
	if providerReg.MatchString(body) {
		log.Println(providerReg.FindStringSubmatch(body)[1])
	}
}
