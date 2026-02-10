package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	required := []string{"DSN", "SALT"}
	missing := []string{}

	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) != 0 {
		log.Fatal("Missing env.vars.: ", missing)
	}
}
