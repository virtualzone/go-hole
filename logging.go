package main

import (
	"log"
	"net"

	"github.com/miekg/dns"
)

var queryTypeNames = map[uint16]string{
	dns.TypeNone:       "None",
	dns.TypeA:          "A",
	dns.TypeNS:         "NS",
	dns.TypeMD:         "MD",
	dns.TypeMF:         "MF",
	dns.TypeCNAME:      "CNAME",
	dns.TypeSOA:        "SOA",
	dns.TypeMB:         "MB",
	dns.TypeMG:         "MG",
	dns.TypeMR:         "MR",
	dns.TypeNULL:       "NULL",
	dns.TypePTR:        "PTR",
	dns.TypeHINFO:      "HINFO",
	dns.TypeMINFO:      "MINFO",
	dns.TypeMX:         "MX",
	dns.TypeTXT:        "TXT",
	dns.TypeRP:         "RP",
	dns.TypeAFSDB:      "AFSDB",
	dns.TypeX25:        "X25",
	dns.TypeISDN:       "ISDN",
	dns.TypeRT:         "RT",
	dns.TypeNSAPPTR:    "NSAPPTR",
	dns.TypeSIG:        "SIG",
	dns.TypeKEY:        "KEY",
	dns.TypePX:         "PX",
	dns.TypeGPOS:       "GPOS",
	dns.TypeAAAA:       "AAAA",
	dns.TypeLOC:        "LOC",
	dns.TypeNXT:        "NXT",
	dns.TypeEID:        "EID",
	dns.TypeNIMLOC:     "NIMLOC",
	dns.TypeSRV:        "SRV",
	dns.TypeATMA:       "ATMA",
	dns.TypeNAPTR:      "NAPTR",
	dns.TypeKX:         "KX",
	dns.TypeCERT:       "CERT",
	dns.TypeDNAME:      "DNAME",
	dns.TypeOPT:        "OPT",
	dns.TypeAPL:        "APL",
	dns.TypeDS:         "DS",
	dns.TypeSSHFP:      "SSHFP",
	dns.TypeRRSIG:      "RRSIG",
	dns.TypeNSEC:       "NSEC",
	dns.TypeDNSKEY:     "DNSKEY",
	dns.TypeDHCID:      "DHCID",
	dns.TypeNSEC3:      "NSEC3",
	dns.TypeNSEC3PARAM: "NSEC3PARAM",
	dns.TypeTLSA:       "TLSA",
	dns.TypeSMIMEA:     "SMIMEA",
	dns.TypeHIP:        "HIP",
	dns.TypeNINFO:      "NINFO",
	dns.TypeRKEY:       "RKEY",
	dns.TypeTALINK:     "TALINK",
	dns.TypeCDS:        "CDS",
	dns.TypeCDNSKEY:    "CDNSKEY",
	dns.TypeOPENPGPKEY: "OPENPGPKEY",
	dns.TypeCSYNC:      "CSYNC",
	dns.TypeZONEMD:     "ZONEMD",
	dns.TypeSVCB:       "SVCB",
	dns.TypeHTTPS:      "HTTPS",
	dns.TypeSPF:        "SPF",
	dns.TypeUINFO:      "UINFO",
	dns.TypeUID:        "UID",
	dns.TypeGID:        "GID",
	dns.TypeUNSPEC:     "UNSPEC",
	dns.TypeNID:        "NID",
	dns.TypeL32:        "L32",
	dns.TypeL64:        "L64",
	dns.TypeLP:         "LP",
	dns.TypeEUI48:      "EUI48",
	dns.TypeEUI64:      "EUI64",
	dns.TypeURI:        "URI",
	dns.TypeCAA:        "CAA",
	dns.TypeAVC:        "AVC",
	dns.TypeTKEY:       "TKEY",
	dns.TypeTSIG:       "TSIG",
	dns.TypeIXFR:       "IXFR",
	dns.TypeAXFR:       "AXFR",
	dns.TypeMAILB:      "MAILB",
	dns.TypeMAILA:      "MAILA",
	dns.TypeANY:        "ANY",
	dns.TypeTA:         "TA",
	dns.TypeDLV:        "DLV",
	dns.TypeReserved:   "Reserved",
}

func getQueryTypeText(qtype uint16) string {
	res := queryTypeNames[qtype]
	if res == "" {
		res = "Unknown"
	}
	return res
}

func logQueryResult(source net.Addr, name string, qtype uint16, result string) {
	log.Printf("Query from %s for %s type %s %s\n", source.String(), name, getQueryTypeText(qtype), result)
}
