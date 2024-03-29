package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

type Configuration struct {
	DatabaseName      string
	DatabaseHost      string
	DatabaseUser      string
	DatabasePassword  string
	MigrateToVersion  string
	MigrationLocation string
	AccessSecret      string
	RefreshSecret     string
	RedisHost         string
	RedisPort         string
}

func GetConfiguration() Configuration {
	migrationLocation, set := os.LookupEnv("MIGRATION_LOCATION")
	if !set {
		migrationLocation = "migrations"
	}

	migrateToVersion, set := os.LookupEnv("MIGRATE")
	if !set {
		migrateToVersion = "latest"
	}
	err := godotenv.Load(filepath.Join(".env"))
	if err != nil {
		log.Println(err)
	}

	return Configuration{
		DatabaseName:      os.Getenv("DB_NAME"),
		DatabaseHost:      os.Getenv("DB_HOST"),
		DatabaseUser:      os.Getenv("DB_USER"),
		DatabasePassword:  os.Getenv("DB_PASSWORD"),
		MigrateToVersion:  migrateToVersion,
		MigrationLocation: migrationLocation,
		AccessSecret:      os.Getenv("ACCESS_SECRET"),
		RefreshSecret:     os.Getenv("REFRESH_SECRET"),
		RedisPort:         os.Getenv("REDIS_PORT"),
		RedisHost:         os.Getenv("REDIS_HOST"),
	}
}
