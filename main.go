package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"reception/app/model"
	"reception/app/server"
	"reception/app/service"
	"reception/app/store/database"
	"strconv"
)

func init() {
	godotenv.Load(".env")
}
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
	config.Debug = os.Getenv("DEBUG") == "true"
	// db config
	config.DbName = os.Getenv("DB_NAME")
	config.DbUser = os.Getenv("DB_USER")
	config.DbHost = os.Getenv("DB_HOST")
	config.DbPort = os.Getenv("DB_PORT")
	config.DbPassword = os.Getenv("DB_PASSWORD")
	// server config
	config.BindHost = os.Getenv("BIND_HOST")
	config.BindPort = os.Getenv("BIND_PORT")
	config.JwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	per, err := strconv.Atoi(os.Getenv("PER"))
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
