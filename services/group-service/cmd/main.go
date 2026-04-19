package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/split/services/group-service/internal/handler"
	"github.com/split/services/group-service/internal/repository"
	"github.com/split/services/group-service/internal/service"
	pb "github.com/split/services/group-service/proto"
)

func main() {
	grpcPort := os.Getenv("PORT")
	if grpcPort == "" {
		grpcPort = getEnv("GRPC_PORT", "50052")
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
		dbName := getEnv("DB_NAME", "group_db")

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

	repo := repository.NewGroupRepository(db)
	svc := service.NewGroupService(repo)
	hdlr := handler.NewGroupHandler(svc)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGroupServiceServer(grpcServer, hdlr)

	log.Printf("Group service listening on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	migration := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS groups (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(255) NOT NULL,
		description TEXT DEFAULT '',
		currency VARCHAR(10) DEFAULT 'USD',
		invite_code VARCHAR(20) UNIQUE,
		owner_id UUID NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS group_members (
		group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
		user_id UUID NOT NULL,
		role VARCHAR(20) NOT NULL DEFAULT 'member',
		joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		PRIMARY KEY (group_id, user_id)
	);

	CREATE INDEX IF NOT EXISTS idx_group_members_user ON group_members(user_id);
	CREATE INDEX IF NOT EXISTS idx_groups_invite_code ON groups(invite_code);
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
