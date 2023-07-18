package tls

import (
	"crypto/tls"
	"fmt"
	"time"
)

type SSLStatus struct {
	Domain    string
	ValidFrom time.Time
	ValidTo   time.Time
	Issuer    string
}

func CheckCertificate(domain string) (*SSLStatus, error) {
	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates

	if len(certs) > 0 {
		cert := certs[0]
		status := &SSLStatus{
			Domain:    domain,
			ValidFrom: cert.NotBefore,
			ValidTo:   cert.NotAfter,
			Issuer:    cert.Issuer.CommonName,
		}
		return status, nil
	}

	return nil, fmt.Errorf("no certificates found")
}
