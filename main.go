package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"store-api/app/misc"
	"store-api/app/model"
	"store-api/app/server"
	"store-api/app/service"
	"store-api/app/store/database"
	"strconv"

	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"
)

func main() {
	config := loadConfig()
	if config.Debug {
		config.DbHost = "localhost"
	}
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)

	db, err := initDB(dbUrl)
	if err != nil {
		fmt.Println("database connection error", err)
	}

	store := database.New(db, config.Per)
	service := service.NewService()
	server := server.NewServer(store, service, config.JwtSecretKey)
	// cors setup
	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Access-Control-Allow-Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST"})

	handlers := handlers.CORS(originsOk, headersOk, methodsOk)(server)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", config.BindHost, config.BindPort), handlers)
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig() model.AppConfig {
	config := model.AppConfig{}
	keys := []string{
		"DEBUG",
		"DB_NAME",
		"DB_USER",
		"DB_HOST",
		"DB_PUBLIC_PORT",
		"DB_PASSWORD",
		"BIND_HOST",
		"BIND_PORT",
		"JWT_SECRET_KEY",
		"PER",
	}
	cnf := misc.GetConfig(keys)
	config.Debug = cnf["DEBUG"] == "true"
	// db config
	config.DbName = cnf["DB_NAME"]
	config.DbUser = cnf["DB_USER"]
	config.DbHost = cnf["DB_HOST"]
	config.DbPort = cnf["DB_PUBLIC_PORT"]
	config.DbPassword = cnf["DB_PASSWORD"]
	// server config
	config.BindHost = cnf["BIND_HOST"]
	config.BindPort = cnf["BIND_PORT"]
	config.JwtSecretKey = cnf["JWT_SECRET_KEY"]
	per, err := strconv.Atoi(cnf["PER"])
	if err != nil {
		config.Per = 10
	} else {
		config.Per = per
	}
	return config
}

func initDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
