package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

var blacklistRecords = map[string]string{}
var whitelistRecords = []string{}

func queryBlacklist(name string, qtype uint16) ([]dns.RR, error) {
	if isWhitelisted(name) {
		return nil, errors.New("record is whitelisted, not checking against blacklist database")
	}
	ip := blacklistRecords[name]
	if ip == "" {
		return nil, errors.New("record not found in blacklist database")
	}
	rr, err := dns.NewRR(fmt.Sprintf("%s A %s", name, ip))
	return []dns.RR{rr}, err
}

func isWhitelisted(name string) bool {
	for _, cur := range whitelistRecords {
		if cur == name {
			return true
		}
	}
	return false
}

func updateBlacklistRecords() {
	log.Println("Updating blacklist database...")
	blacklistRecords = make(map[string]string, 0)
	for _, url := range GetConfig().BlacklistSources {
		processBlacklistSource(url)
	}
}

func updateWhitelistRecords() {
	log.Println("Updating whitelist database...")
	whitelistRecords = make([]string, 0)
	for _, name := range GetConfig().Whitelist {
		whitelistRecords = append(whitelistRecords, strings.ToLower(strings.TrimSpace(name))+".")
	}
}

func processBlacklistSource(url string) error {
	data, err := getUrlData(url)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(data)
	fileScanner := bufio.NewScanner(reader)
	fileScanner.Split(bufio.ScanLines)
	re := regexp.MustCompile(`\s+`)
	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		if line != "" && line[0] != '#' {
			split := re.Split(line, -1)
			if len(split) == 2 {
				blacklistRecords[strings.ToLower(split[1])+"."] = split[0]
			}
		}
	}
	return nil
}

func getUrlData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("received invalid http response code " + strconv.Itoa(resp.StatusCode) + "for url " + url)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
