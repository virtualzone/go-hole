package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var blacklistRecords = []string{}
var whitelistRecords = []string{}

func queryBlacklist(name string, qtype uint16) ([]dns.RR, error) {
	if isWhitelisted(name) {
		return nil, errors.New("record is whitelisted, not checking against blacklist database")
	}
	if !isBlacklisted(name) {
		return nil, errors.New("record not found in blacklist database")
	}
	return []dns.RR{}, nil
}

func isBlacklisted(name string) bool {
	if GetConfig().BlacklistEverything {
		return true
	}
	for _, cur := range blacklistRecords {
		if cur == name {
			return true
		}
	}
	return false
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
	list := make([]string, 0)
	for _, url := range GetConfig().BlacklistSources {
		processBlacklistSource(url, &list)
	}
	blacklistRecords = list
	log.Printf("Blacklist database updated, %d records\n", len(blacklistRecords))
}

func initBlacklistRenewal() {
	if GetConfig().BlacklistRenewal < 1 {
		return
	}
	ticker := time.NewTicker(time.Minute * time.Duration(GetConfig().BlacklistRenewal))
	go func() {
		for {
			<-ticker.C
			updateBlacklistRecords()
		}
	}()
}

func updateWhitelistRecords() {
	log.Println("Updating whitelist database...")
	whitelistRecords = make([]string, 0)
	for _, name := range GetConfig().Whitelist {
		whitelistRecords = append(whitelistRecords, strings.ToLower(strings.TrimSpace(name))+".")
	}
}

func processBlacklistSource(url string, list *[]string) error {
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
			if isValidBlacklistSourceRecord(split) {
				if len(split) == 2 {
					*list = append(*list, strings.ToLower(split[1])+".")
				} else if len(split) == 1 {
					*list = append(*list, strings.ToLower(split[0])+".")
				}
			}
		}
	}
	return nil
}

func isValidBlacklistSourceRecord(split []string) bool {
	if len(split) == 0 {
		return false
	}
	if len(split) > 2 {
		return false
	}
	if len(split) == 2 {
		if split[0] != "0.0.0.0" {
			return false
		}
		if split[0] == "0.0.0.0" && split[1] == "0.0.0.0" {
			return false
		}
	}
	return true
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
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
