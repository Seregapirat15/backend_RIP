// Run: go run scripts/gen_password.go
package main

import (
	"fmt"
	"lab4/internal/auth"
)

func main() {
	for _, pwd := range []string{"modpass123", "password"} {
		h, err := auth.HashPassword(pwd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("password=%q -> %s\n", pwd, h)
	}
}
