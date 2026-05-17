package main

import (
	"fmt"

	"tool-go/internal/library/password"
)

func main() {
	passwords := []string{"admin123"}

	for _, pwd := range passwords {
		hash, salt, err := password.CreatePassword(pwd)
		if err != nil {
			fmt.Printf("Error generating password %s: %v\n", pwd, err)
			continue
		}
		fmt.Printf("Password: %s\n", pwd)
		fmt.Printf("Salt: %s\n", salt)
		fmt.Printf("Hash: %s\n", hash)
		fmt.Printf("SQL: INSERT INTO \"user\" (\"username\", \"password\", \"salt\", \"nickname\", \"status\") VALUES ('admin', '%s', '%s', '超级管理员', 1);\n\n", hash, salt)
	}
}
