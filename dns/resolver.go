package dns

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
	
	
)

func Resolve(nameserver string, rdtype []string, fqdn string) (*dns.Msg, time.Duration) {
	var in *dns.Msg
	var rtt time.Duration
	var err error
	for _, recordType := range rdtype {
		m1 := new(dns.Msg)
		m1.Id = dns.Id()
		m1.RecursionDesired = true
		m1.Question = make([]dns.Question, 1)
		if qType, ok := dns.StringToType[recordType]; ok {
			fqdn = dns.Fqdn(fqdn)
			m1.Question[0] = dns.Question{
				Name: fqdn, 
				Qtype: qType, 
				Qclass: dns.ClassINET,
			}
		} else {
			fmt.Println("Unsupported RD type:", rdtype)
		}
		c := new(dns.Client)
		c.Net = "udp"
		in, rtt, err = c.Exchange(m1, fmt.Sprintf("%s:53", nameserver))
		if err != nil {
			fmt.Println(err)
		}
		FormatDNSResponse(in, rtt)
		
	}
	return nil,0
}
