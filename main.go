package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mhmdfathan/funsies-be/config"
)

func main() {
	httpServer := http.NewServeMux()
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	  return
	}
	config.DatabaseInit()

	fmt.Println("Memoraire backend is listening at port:", os.Getenv("APP_PORT"))
	http.ListenAndServe(":" + os.Getenv("APP_PORT"), httpServer)
}