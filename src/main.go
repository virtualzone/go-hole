package main

import "log"

func main() {
	log.Printf("Starting Go-hole %s...\n", AppVersion)
	GetConfig().ReadConfig()
	GetConfig().Print()
	initServer()
	listenAndServe()
}

func initServer() {
	updateLocalRecords()
	updateBlacklistRecords()
	updateWhitelistRecords()
}
