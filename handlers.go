package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func addDomain(c *fiber.Ctx) error {
	log.Println(c.Body())
	// Parse the request body
	var domain Domain
	if err := c.BodyParser(&domain); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	log.Println("Adding domain: ", domain.Name)

	// Add the domain to the database
	if err := db.AddDomain(domain.Name); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot add domain",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Domain added",
	})
}

func getDomains(c *fiber.Ctx) error {
	// Get the domains from the database
	domains, err := db.GetDomains()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot get domains",
		})
	}

	// Respond with the domains
	return c.JSON(domains)
}
