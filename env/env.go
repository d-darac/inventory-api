package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	DB_URL     string
	HOST       string
	IV         string
	MASTER_KEY string
	PLATFORM   string
	PORT       string
}

func GetEnv() env {
	platform, ok := os.LookupEnv("PLATFORM")
	if !ok || platform == "dev" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("couldn't load env variables: %v", err)
			os.Exit(1)
		}
		platform = "dev"
	}

	dbUrl, ok := os.LookupEnv("DB_URL")
	if !ok {
		log.Fatalln("env variable 'DB_URL' not set")
		os.Exit(1)
	}
	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatalln("env variable 'HOST' not set")
		os.Exit(1)
	}
	iv, ok := os.LookupEnv("IV")
	if !ok {
		log.Fatalln("env variable 'IV' not set")
		os.Exit(1)
	}
	key, ok := os.LookupEnv("MASTER_KEY")
	if !ok {
		log.Fatalln("env variable 'MASTER_KEY' not set")
		os.Exit(1)
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatalln("env variable 'PORT' not set")
		os.Exit(1)
	}

	return env{
		DB_URL:     dbUrl,
		HOST:       host,
		IV:         iv,
		MASTER_KEY: key,
		PLATFORM:   platform,
		PORT:       port,
	}
}
