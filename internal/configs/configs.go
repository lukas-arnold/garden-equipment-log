package configs

import (
	"io/fs"
	"os"

	"github.com/lukas-arnold/garden-equipment-log/web"
)

const (
	defaultPort        = ":8080"
	defaultStorageFile = "data/strength-tracker.json"
	defaultLanguage    = "en"
)

func GetPort() string {
	return getEnv("PORT", defaultPort)
}

func GetStorageFile() string {
	return getEnv("STORAGE_FILE", defaultStorageFile)
}

func GetLanguage() string {
	return getEnv("LANGUAGE", defaultLanguage)
}

func GetWebFiles() fs.FS {
	return web.WebFiles
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
