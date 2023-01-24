package config

import (
	"os"
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

	return Configuration{
		DatabaseName:      "db_nix_junior",
		DatabaseHost:      "localhost:8081",
		DatabaseUser:      "admin",
		DatabasePassword:  "password",
		MigrateToVersion:  migrateToVersion,
		MigrationLocation: migrationLocation,
		AccessSecret:      "access",
		RefreshSecret:     "refresh",
		RedisPort:         "6379",
		RedisHost:         "localhost",
	}
}
