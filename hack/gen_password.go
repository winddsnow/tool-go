package main

import (
	"fmt"

	"tool-go/internal/library/password"
)

func main() {
	passwords := []string{"walter"}

	for _, pwd := range passwords {
		hash, salt, err := password.CreatePassword(pwd)
		if err != nil {
			fmt.Printf("Error generating password %s: %v\n", pwd, err)
			continue
		}
		fmt.Printf("Password: %s\n", pwd)
		fmt.Printf("Salt: %s\n", salt)
		fmt.Printf("Hash: %s\n", hash)
		fmt.Printf("SQL: INSERT INTO \"user\" (\"username\", \"password\", \"salt\", \"nickname\", \"status\") VALUES ('walter', '%s', '%s', '本地开发', 1);\n\n", hash, salt)
	}
}
