package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Hash from insert_user.sql
	h1 := []byte("$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/X4aMqNw6FKP6PZQ7i")
	// Hash from update_user.sql
	h2 := []byte("$2b$12$AZqZ1LohtF9/gSLKylHbNeZS5KjpMAZ09FhN7j9II/V4H0E9aZ7o6")

	passwords := []string{"admin", "admin123"}
	hashes := [][]byte{h1, h2}
	hashLabels := []string{"insert_user.sql hash", "update_user.sql hash"}

	for _, p := range passwords {
		fmt.Printf("\nPassword: %s\n", p)
		for i, h := range hashes {
			err := bcrypt.CompareHashAndPassword(h, []byte(p))
			if err == nil {
				fmt.Printf("  ✓ MATCHES %s\n", hashLabels[i])
			} else {
				fmt.Printf("  ✗ NO match %s\n", hashLabels[i])
			}
		}
	}
}
