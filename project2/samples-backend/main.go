package main

import (
	"fmt"
	"log"
	"net/http"

	db "sample-backend-go/internal/infrastructure/database"
	"sample-backend-go/internal/infrastructure/middleware"
	mergeRouter "sample-backend-go/internal/interface/routers"

	"github.com/joho/godotenv"
)

func init() {
	// .envファイルが存在する場合は読み込む
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	// Connect to Cassandra Cluster
	cSession, err := db.NewSession("cassandra-node1:9042", "sample_dev")
	if err != nil {
		panic(err)
	}
	defer cSession.Close()

	// Connect to PostgreSQL
	pgPool, err := db.NewPostgresSession("user=postgres password=secret host=db-postgres dbname=sample_dev port=5432 sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer pgPool.Close()

	// Initialize routers
	RouterBundled := mergeRouter.SetupMergedRouter(pgPool, cSession)

	RouterWithCORS := middleware.SetupCORS(RouterBundled)

	// Set up HTTP server with separate route prefixes
	http.Handle("/api/", http.StripPrefix("/api", RouterWithCORS))

	// Start the server
	fmt.Println("Server is starting on :5000")
	http.ListenAndServe(":5000", nil)
}
