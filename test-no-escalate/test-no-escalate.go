package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	fmt.Printf("Current UID: %d, GID: %d\n", os.Getuid(), os.Getgid())
	fmt.Printf("Current EUID: %d, EGID: %d\n", os.Geteuid(), os.Getegid())
	
	// Test that both seteuid(0) and setegid(0) fail as expected
	euidError := syscall.Seteuid(0)
	egidError := syscall.Setegid(0)
	
	if euidError != nil && egidError != nil {
		fmt.Printf("Got expected error when setting EUID to 0: %v\n", euidError)
		fmt.Printf("Got expected error when setting EGID to 0: %v\n", egidError)
		// This is the expected behavior - exit with success
		os.Exit(0)
	} else {
		// At least one of them succeeded, which is a security vulnerability
		if euidError == nil {
			fmt.Printf("ERROR: Successfully set EUID to 0. New EUID: %d\n", os.Geteuid())
		}
		if egidError == nil {
			fmt.Printf("ERROR: Successfully set EGID to 0. New EGID: %d\n", os.Getegid())
		}
		// Exit with failure
		os.Exit(1)
	}
}
