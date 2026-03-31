package main

import (
	"flag"
	"fmt"
	"strings"
)

// Custom flag type which can be used to collect a range of different record types to test.
var recordTypeFlag string

func main() {
	var nsFlag = flag.String("ns", "", "DNS nameserver address [IP Address or FQDN]")
	var dnsRecordFlag = flag.String("dns", "www.example.com", "DNS record to test")
	var recordTypeFlag = flag.String("rdtype", "A, AAAA, CNAME, ", "Record type [A or AAAA], this option accepts a range of different values. For multiple values please pass a comma separate list")
	flag.Parse()

	fmt.Printf("ns=%s\n", *nsFlag)
	fmt.Printf("dns=%s\n", *dnsRecordFlag)
	fmt.Printf("rdtype=%s\n", *recordTypeFlag)
	var recordTypeList []string
	// Split the list of rdtype records into a list of types which can be used more efficiently.
	rawList := strings.Split(*recordTypeFlag, ",")
	for _, item := range rawList {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			recordTypeList = append(recordTypeList, trimmed)
		}
	}
	for _, recordType := range recordTypeList {
		fmt.Printf("%s\n", recordType)
	}
}
