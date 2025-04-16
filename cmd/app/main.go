package main

import (
	"database/sql"
	"fmt"
	"gateway/internal/repository"
	"gateway/internal/service"
	"gateway/internal/web/server"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getEnv(key, defaultValue string) string{
	if value := os.Getenv(key); value != ""{
		return value
	}
	return defaultValue
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro loading .env file")
	}
	//string de conexao
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "gateway"),
		getEnv("DB_SSL_MODE", "disable"),
	)
	//conexao com db
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro connecting to database: ", err)
	}
	defer db.Close()

	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)

	invoiceRepository := repository.NewInvoiceRepository(db)
	invoiceService := service.NewInvoiceService(invoiceRepository, *accountService)

	port := getEnv("HTTP_PORT", "8080")
	srv := server.NewServer(accountService, invoiceService, port)
	srv.ConfigureRoutes()
	log.Println("Starting server on port", port)

	srv.Start()
	if err != nil{
		log.Fatal("Error starting server: ", err)
	}
}