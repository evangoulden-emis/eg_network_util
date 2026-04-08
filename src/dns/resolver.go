package dns

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

func Resolve(nameserver string, rdtype []string, fqdn string, useRaw bool) (*dns.Msg, time.Duration) {
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
		
		if useRaw {
			fmt.Printf("Using raw DNS query to %s for %s %s\n", nameserver, fqdn, recordType)
			// Use raw sockets to bypass interception
			in, rtt, err = resolveRawUDP(nameserver, m1)
			if err != nil {
				fmt.Printf("Raw UDP failed, trying TCP: %v\n", err)
				in, rtt, err = resolveRawTCP(nameserver, m1)
				if err != nil {
					fmt.Printf("Raw TCP also failed: %v\n", err)
					continue
				}
			}
		} else {
			// Use standard DNS client
			c := new(dns.Client)
			c.Net = "udp"
			in, rtt, err = c.Exchange(m1, fmt.Sprintf("%s:53", nameserver))
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		
		FormatDNSResponse(in, rtt)
	}
	return nil, 0
}

// resolveRawUDP sends DNS query using raw UDP socket to bypass system DNS interception
func resolveRawUDP(nameserver string, query *dns.Msg) (*dns.Msg, time.Duration, error) {
	fmt.Printf("Attempting raw UDP connection to %s:53\n", nameserver)
	
	// Pack the DNS query message
	queryData, err := query.Pack()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to pack query: %v", err)
	}
	fmt.Printf("Packed query: %d bytes\n", len(queryData))

	// Create UDP connection
	conn, err := net.Dial("udp", fmt.Sprintf("%s:53", nameserver))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to connect to %s: %v", nameserver, err)
	}
	defer conn.Close()
	fmt.Printf("UDP connection established\n")

	// Set timeout
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Send query
	start := time.Now()
	n, err := conn.Write(queryData)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send query: %v", err)
	}
	fmt.Printf("Sent %d bytes\n", n)

	// Read response
	responseData := make([]byte, 4096)
	n, err = conn.Read(responseData)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response: %v", err)
	}
	rtt := time.Since(start)
	fmt.Printf("Received %d bytes in %v\n", n, rtt)

	// Unpack the DNS response
	response := new(dns.Msg)
	err = response.Unpack(responseData[:n])
	if err != nil {
		return nil, 0, fmt.Errorf("failed to unpack response: %v", err)
	}

	return response, rtt, nil
}

// resolveRawTCP sends DNS query using raw TCP socket (less likely to be intercepted)
func resolveRawTCP(nameserver string, query *dns.Msg) (*dns.Msg, time.Duration, error) {
	// Pack the DNS query message
	queryData, err := query.Pack()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to pack query: %v", err)
	}

	// DNS over TCP requires a 2-byte length prefix
	lengthPrefix := make([]byte, 2)
	lengthPrefix[0] = byte(len(queryData) >> 8)
	lengthPrefix[1] = byte(len(queryData))

	// Create TCP connection
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:53", nameserver))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to connect to %s: %v", nameserver, err)
	}
	defer conn.Close()

	// Set timeout
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Send length prefix + query
	start := time.Now()
	_, err = conn.Write(append(lengthPrefix, queryData...))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send query: %v", err)
	}

	// Read length prefix
	_, err = conn.Read(lengthPrefix)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read length prefix: %v", err)
	}

	// Calculate response length
	responseLength := int(lengthPrefix[0])<<8 | int(lengthPrefix[1])
	responseData := make([]byte, responseLength)

	// Read response
	_, err = conn.Read(responseData)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response: %v", err)
	}
	rtt := time.Since(start)

	// Unpack the DNS response
	response := new(dns.Msg)
	err = response.Unpack(responseData)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to unpack response: %v", err)
	}

	return response, rtt, nil
}
