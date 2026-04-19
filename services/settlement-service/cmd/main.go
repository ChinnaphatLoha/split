package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/split/services/settlement-service/internal/handler"
	"github.com/split/services/settlement-service/internal/repository"
	"github.com/split/services/settlement-service/internal/service"
	pb "github.com/split/services/settlement-service/proto"
)

func main() {
	grpcPort := os.Getenv("PORT")
	if grpcPort == "" {
		grpcPort = getEnv("GRPC_PORT", "50054")
	}

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
		dbName := getEnv("DB_NAME", "settlement_db")

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

	if err := runMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	repo := repository.NewSettlementRepository(db)
	svc := service.NewSettlementService(repo)
	hdlr := handler.NewSettlementHandler(svc)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSettlementServiceServer(grpcServer, hdlr)

	log.Printf("Settlement service listening on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	migration := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS settlements (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		group_id UUID NOT NULL,
		from_user_id UUID NOT NULL,
		to_user_id UUID NOT NULL,
		amount DECIMAL(12,2) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_settlements_group ON settlements(group_id);
	CREATE INDEX IF NOT EXISTS idx_settlements_from_user ON settlements(from_user_id);
	CREATE INDEX IF NOT EXISTS idx_settlements_to_user ON settlements(to_user_id);
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
