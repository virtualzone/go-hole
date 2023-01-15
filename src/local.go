package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/miekg/dns"
)

type LocalRecordTarget struct {
	Target string
	Qtype  uint16
}

var localRecords = map[string][]LocalRecordTarget{}

func queryLocal(name string, qtype uint16) ([]dns.RR, error) {
	target, ok := localRecords[name]
	if !ok {
		return nil, errors.New("record not found in local database")
	}
	res := make([]dns.RR, 0)
	for _, record := range target {
		if record.Qtype == qtype {
			rr, err := dns.NewRR(fmt.Sprintf("%s %s %s", name, getQueryTypeText(record.Qtype), record.Target))
			if err != nil {
				log.Println(err)
				return []dns.RR{}, err
			}
			res = append(res, rr)
		}
	}
	if len(res) == 0 {
		return nil, errors.New("no record for requested query type found in local database")
	}
	return res, nil
}

func updateLocalRecords() {
	log.Println("Updating local address database...")
	localRecords = make(map[string][]LocalRecordTarget, 0)
	for _, item := range GetConfig().LocalAddresses {
		records := make([]LocalRecordTarget, 0)
		for _, target := range item.Target {
			qtype, err := getQueryTypeUint(target.Type)
			if err != nil {
				log.Printf("Ignoring unknown record type %s for %s\n", target.Type, item.Name)
				continue
			}
			record := LocalRecordTarget{
				Target: target.Address,
				Qtype:  qtype,
			}
			records = append(records, record)
		}
		localRecords[strings.ToLower(item.Name)+"."] = records
	}
}
