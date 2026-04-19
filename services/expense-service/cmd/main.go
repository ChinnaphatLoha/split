package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/split/services/expense-service/internal/handler"
	"github.com/split/services/expense-service/internal/repository"
	"github.com/split/services/expense-service/internal/service"
	pb "github.com/split/services/expense-service/proto"
)

func main() {
	grpcPort := os.Getenv("PORT")
	if grpcPort == "" {
		grpcPort = getEnv("GRPC_PORT", "50053")
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
		dbName := getEnv("DB_NAME", "expense_db")

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

	repo := repository.NewExpenseRepository(db)
	svc := service.NewExpenseService(repo)
	hdlr := handler.NewExpenseHandler(svc)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExpenseServiceServer(grpcServer, hdlr)

	log.Printf("Expense service listening on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	migration := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS expenses (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		group_id UUID NOT NULL,
		payer_id UUID NOT NULL,
		amount DECIMAL(12,2) NOT NULL,
		description TEXT DEFAULT '',
		split_type VARCHAR(20) NOT NULL DEFAULT 'EQUAL',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS expense_splits (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		expense_id UUID NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
		user_id UUID NOT NULL,
		amount DECIMAL(12,2) NOT NULL,
		percentage DECIMAL(5,2) DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_expenses_group ON expenses(group_id);
	CREATE INDEX IF NOT EXISTS idx_expenses_payer ON expenses(payer_id);
	CREATE INDEX IF NOT EXISTS idx_expense_splits_expense ON expense_splits(expense_id);
	CREATE INDEX IF NOT EXISTS idx_expense_splits_user ON expense_splits(user_id);
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
