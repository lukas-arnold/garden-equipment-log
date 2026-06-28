package main

import (
	"log"
	"net/http"

	"github.com/lukas-arnold/garden-equipment-log/internal/configs"
	"github.com/lukas-arnold/garden-equipment-log/internal/handler"
	"github.com/lukas-arnold/garden-equipment-log/internal/language"
	"github.com/lukas-arnold/garden-equipment-log/internal/storage"
)

func main() {
	err := language.LoadLanguages()

	if err != nil {
		log.Fatal(err)
	}

	store := storage.New(
		configs.GetStorageFile(),
	)

	h := handler.New(
		store,
	)

	server := createServer(
		h,
	)

	log.Printf(
		"Garden Equipment Log running on %s",
		configs.GetPort(),
	)

	log.Fatal(
		http.ListenAndServe(
			configs.GetPort(),
			server,
		),
	)
}
