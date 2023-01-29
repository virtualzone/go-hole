package main

import "log"

func main() {
	log.Printf("Starting Go-hole %s...\n", AppVersion)
	GetConfig().ReadConfig()
	GetConfig().Print()
	initServer()
	initBlacklistRenewal()
	listenAndServe()
}

func initServer() {
	initLogging()
	GetUpstreamCache().Init()
	updateLocalRecords()
	updateBlacklistRecords()
	updateWhitelistRecords()
}
