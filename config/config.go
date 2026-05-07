package config

type Config struct {
	DatabaseURL string
	SlotName	string
	PublicationName string
}

func DefaultConfig() Config {
	return Config{
		DatabaseURL: "postgres://user:password@localhost:5432/dbname",
		SlotName: "pgstream_slot",
		PublicationName: "pgstream_pub",
	}
}