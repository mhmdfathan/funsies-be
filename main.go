package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mhmdfathan/funsies-be/config"
	"github.com/mhmdfathan/funsies-be/routes"
	"github.com/mhmdfathan/funsies-be/utils"
)

func main() {


	httpServer := http.NewServeMux()
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	  return
	}

	//database
	config.DatabaseInit()

	//cron for checking pending users
	utils.StartCron(config.DB)

	//routes
	routes.UserRoutes(httpServer)

	fmt.Println("Memoraire backend is listening at port:", os.Getenv("APP_PORT"))
	http.ListenAndServe(":" + os.Getenv("APP_PORT"), httpServer)
}