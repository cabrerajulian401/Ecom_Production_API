package env

import "os"

/* Want to create a function to get environment variables
because in case env varible does not exist it returns a fall back
value */

func GetString(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}
