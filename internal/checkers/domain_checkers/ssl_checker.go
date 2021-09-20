package domain_checkers

import (
	"crypto/tls"
	"fmt"
	"time"
)

const Day = time.Hour * 24

type SSLChecker struct {
	domains []string
}

func NewSSLChecker(domains []string) *SSLChecker {
	return &SSLChecker{
		domains: domains,
	}
}

func (c *SSLChecker) Check() bool {
	for _, domain := range c.domains {
		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", domain), nil)
		if err != nil {
			// fmt.Printf("ssl error: %s\n", domain)
			return false
		}

		if err = conn.VerifyHostname(domain); err != nil {
			// fmt.Printf("hostname ssl error: %s\n", domain)
			return false
		}
		expiry := conn.ConnectionState().PeerCertificates[0].NotAfter

		if time.Now().Add(time.Duration(14) * Day).After(expiry) {
			fmt.Printf("domain: %s\nExpiry: %v\n", domain, expiry.Format(time.RFC850))
			return false
		}
	}

	return true
}
