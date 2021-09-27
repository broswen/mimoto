package main

import (
	"log"

	"github.com/broswen/mimoto/internal/server"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("load dotenv: %v", err)
	}

	server, err := server.New()
	if err != nil {
		log.Fatalf("init server: %v", err)
	}

	server.Routes()
	if err != nil {
		log.Fatalf("server routes: %v", err)
	}

	if err := server.Listen(); err != nil {
		log.Fatalf("listen server: %v", err)
	}

}
