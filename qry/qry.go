package qry

import (
	"crypto/tls"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/google/goterm/term"
	"github.com/miekg/dns"
)

func getRandType() string {
	rand.Seed(time.Now().UnixNano())
	types := []string{"A", "SOA", "MX", "TXT", "AAAA"}
	return types[rand.Intn(len(types))]
}

// Send a single DNS query
func SimpleQuery(server string,
	port string,
	qname string,
	qtype string,
	responses chan Response,
	proto string,
	wg *sync.WaitGroup,
	noverify bool) {

	s_server := net.JoinHostPort(server, port)
	typ := getRandType()
	qrytype := Qtype(typ)

	question := new(dns.Msg)
	question.SetQuestion(dns.Fqdn(qname), qrytype)
	c := new(dns.Client)
	c.Dialer = &net.Dialer{
		Timeout: 15 * time.Second,
	}
	c.Timeout = 15 * time.Second
	c.DialTimeout = 15 * time.Second
	c.Net = proto
	if proto == "tcp-tls" {
		var tlc tls.Config
		tlc.InsecureSkipVerify = noverify
		c.TLSConfig = &tlc
	}

	ans, rtt, err := c.Exchange(question, s_server)
	if err != nil {
		log.Println(term.Redf(err.Error()))
	} else {
		var R Response
		R.Rcode = Rcode(ans.Rcode)
		R.Rtt = rtt
		R.Qname = qname
		R.Server = s_server
		R.Qtype = qtype
		responses <- R
	}
	wg.Done()
}
