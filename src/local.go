package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/miekg/dns"
)

var localRecords = map[string]string{}

func queryLocal(name string, qtype uint16) ([]dns.RR, error) {
	ip := localRecords[name]
	if ip == "" {
		return nil, errors.New("record not found in local database")
	}
	rr, err := dns.NewRR(fmt.Sprintf("%s A %s", name, ip))
	return []dns.RR{rr}, err
}

func updateLocalRecords() {
	log.Println("Updating local address database...")
	localRecords = make(map[string]string, 0)
	for _, item := range GetConfig().LocalAddresses {
		localRecords[strings.ToLower(item.Name)+"."] = item.Address
	}
}
