package main

import (
	"log"

	"github.com/robfig/cron/v3"
)

func StartCheckCertificateValidityCronJob(cronPattern string) {
	if cronPattern == "" {
		log.Println("No cron pattern provided, not starting cron job")
		return
	}

	if cronPattern == "-" {
		log.Println("Checking certificate validity once")
		checkCertificateValidity()
		log.Println("Done checking certificate validity")
		return
	}

	c := cron.New(cron.WithSeconds()) // cron with second-level precision
	_, err := c.AddFunc(cronPattern, func() {
		checkCertificateValidity()
	})
	if err != nil {
		log.Fatalf("Cannot create cron job: %v", err)
	}

	log.Println("Starting cron job")
	c.Start()
}

func checkCertificateValidity() {
	// Get all domains from the database
	domains, err := db.GetDomains()
	if err != nil {
		log.Printf("Cannot get domains: %v", err)
		return
	}

	// Check the SSL certificate of each domain
	for _, domain := range domains {
		log.Println("Checking certificate of", domain.Name)
		status, err := CheckCertificate(domain.Name)
		if err != nil {
			log.Printf("Cannot check certificate of %s: %v", domain.Name, err)
			domain.CertificateExpiry.Valid = false
			domain.Status = StatusUnknown
			if err := db.UpdateDomain(domain); err != nil {
				log.Printf("Cannot update domain: %v", err)
			}
			continue
		}

		log.Println(status)

		// Update the domain in the database
		domain.CertificateExpiry.Time = status.ValidTo
		domain.CertificateExpiry.Valid = true
		domain.Status = StatusValid
		domain.Issuer.String = status.Issuer
		domain.Issuer.Valid = true
		if err := db.UpdateDomain(domain); err != nil {
			log.Printf("Cannot update domain: %v", err)
		}
	}
}
