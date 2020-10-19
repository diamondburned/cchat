// Package adapter provides functions to convert typical cchat calls to
// RPC-compatible payloads.
package adapter

import (
	"log"

	"github.com/diamondburned/cchat"
	"github.com/diamondburned/cchat/services"
)

// Main is a helper function that calls Start() wtih the services from
// services.Get().
func Main() error {
	s, errs := services.Get()
	if errs != nil {
		for _, err := range errs {
			if err != nil {
				log.Println(err)
			}
		}
	}

	return Start(s)
}

func Start(services []cchat.Service) error {
	panic("Implement me")
}
