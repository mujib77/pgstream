package config

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL     string
	SlotName        string
	PublicationName string
}

func DefaultConfig() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("no .env file found, using environment variables")
	}

    return Config{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		SlotName:        os.Getenv("SLOT_NAME"),
		PublicationName: os.Getenv("PUBLICATION_NAME"),
	}
}
