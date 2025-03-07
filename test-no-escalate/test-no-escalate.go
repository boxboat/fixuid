package main

import (
	"log"
	"os"
	"syscall"
)

var logger = log.New(os.Stderr, "", 0)

func main() {
	logger.SetPrefix("test-no-escalate: ")
	
	logger.Printf("Current UID: %d, GID: %d", os.Getuid(), os.Getgid())
	logger.Printf("Current EUID: %d, EGID: %d", os.Geteuid(), os.Getegid())
	
	// Test that both seteuid(0) and setegid(0) fail as expected
	euidError := syscall.Seteuid(0)
	egidError := syscall.Setegid(0)
	
	if euidError != nil && egidError != nil {
		logger.Printf("Got expected error when setting EUID to 0: %v", euidError)
		logger.Printf("Got expected error when setting EGID to 0: %v", egidError)
		// This is the expected behavior - exit with success
		os.Exit(0)
	} else {
		// At least one of them succeeded, which is a security vulnerability
		if euidError == nil {
			logger.Printf("ERROR: Successfully set EUID to 0. New EUID: %d", os.Geteuid())
		}
		if egidError == nil {
			logger.Printf("ERROR: Successfully set EGID to 0. New EGID: %d", os.Getegid())
		}
		// Exit with failure
		os.Exit(1)
	}
}
