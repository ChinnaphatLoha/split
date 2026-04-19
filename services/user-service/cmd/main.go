package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/split/services/user-service/internal/handler"
	"github.com/split/services/user-service/internal/repository"
	"github.com/split/services/user-service/internal/service"
	pb "github.com/split/services/user-service/proto"
)

func main() {
	// Database connection
	grpcPort := os.Getenv("PORT")
	if grpcPort == "" {
		grpcPort = getEnv("GRPC_PORT", "50051")
	}
	jwtSecret := getEnv("JWT_SECRET", "split-jwt-secret-key-change-in-production")

	dbUrl := os.Getenv("DATABASE_URL")
	var db *sql.DB
	var err error

	if dbUrl != "" {
		db, err = sql.Open("postgres", dbUrl)
	} else {
		dbHost := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "5432")
		dbUser := getEnv("DB_USER", "split")
		dbPassword := getEnv("DB_PASSWORD", "split_secret")
		dbName := getEnv("DB_NAME", "user_db")

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)
	
		db, err = sql.Open("postgres", dsn)
	}
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Initialize layers
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, jwtSecret)
	hdlr := handler.NewUserHandler(svc)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, hdlr)

	log.Printf("User service listening on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	migration := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		avatar_url TEXT DEFAULT '',
		currency VARCHAR(10) DEFAULT 'USD',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`
	_, err := db.Exec(migration)
	return err
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
