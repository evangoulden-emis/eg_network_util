package dns

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/miekg/dns"
	"github.com/olekukonko/tablewriter"
)

func FormatDNSResponse(msg *dns.Msg, rtt time.Duration) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Index","Name", "Type", "Class", "TTL", "Data"})
	
	for idx, answer := range msg.Answer {
		index := idx + 1
		name := answer.Header().Name
		rtype := dns.TypeToString[answer.Header().Rrtype]
		class := dns.ClassToString[answer.Header().Class]
		ttl := answer.Header().Ttl
		data := answer.String()
		
		table.Append([]string{strconv.Itoa(index), name, rtype, class, strconv.FormatUint(uint64(ttl), 10), data})
	}
	fmt.Printf("Response Time: %v\n", rtt)
	table.Render()
}