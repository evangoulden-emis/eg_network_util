package main

import (
	"flag"
	"fmt"
	"strings"

	"eg_network_util/dns"
)

// Custom flag type which can be used to collect a range of different record types to test.
var recordTypeFlag string

func main() {
	// Collect Flags from the user.
	var nsFlag = flag.String("ns", "8.8.8.8", "DNS nameserver address [IP Address or FQDN]")
	var dnsRecordFlag = flag.String("dns", "www.google.com", "DNS record to test")
	var dnsRecordTypeFlag = flag.String("rdtype", "A, AAAA, CNAME, ", "Record type [A or AAAA], this option accepts a range of different values. For multiple values please pass a comma separate list")
	var _ = flag.Int("asn", 1, "ASN number")                                  // TODO: Add variable name.
	var _ = flag.String("bgp-cidr", "185.13.72.0/23", "CIDR range to lookup") // TODO: Add variable name.
	// Parse all of the flags passed to the application.
	flag.Parse()

	// Validate and resolve any formatting issues passed to the application pertaining to the rdtype values.
	var recordTypeList []string
	// Split the list of rdtype records into a list of types which can be used more efficiently.
	// Rdtype is a DNS lookup type such as A, CNAME, NS etc.
	rawList := strings.Split(*dnsRecordTypeFlag, ",")
	for _, item := range rawList {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			recordTypeList = append(recordTypeList, trimmed)
		}
	}
	for _, recordType := range recordTypeList {
		fmt.Printf("%s\n", recordType)
	}
	dns.Resolve(*nsFlag, recordTypeList, *dnsRecordFlag)
}
