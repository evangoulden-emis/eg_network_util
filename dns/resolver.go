package dns

import (
	"fmt"
	"github.com/miekg/dns"
)

func Resolve(nameserver string, rdtype []string, fqdn string) {
	// Testing merge to main
	for _, recordType := range rdtype {
		m1 := new(dns.Msg)
		m1.Id = dns.Id()
		m1.RecursionDesired = true
		m1.Question = make([]dns.Question, 1)
		if qType, ok := dns.StringToType[recordType]; ok {
			fqdn = dns.Fqdn(fqdn)
			m1.Question[0] = dns.Question{fqdn, qType, dns.ClassINET}
		} else {
			fmt.Println("Unsupported RD type:", rdtype)
		}
		c := new(dns.Client)
		c.Net = "udp"
		in, rtt, err := c.Exchange(m1, fmt.Sprintf("%s:53", nameserver))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Total RTT: ", rtt)
		fmt.Printf("%+v\n", in)
	}
}
