.PHONY: proto-gen build-services dev down clean

# ─── Proto Generation ─────────────────────────────────────────
proto-gen:
	@echo "Generating Go protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user/user.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/group/group.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/expense/expense.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/settlement/settlement.proto
	@echo "Proto generation complete!"

# ─── Docker ───────────────────────────────────────────────────
dev:
	docker-compose up --build

down:
	docker-compose down

clean:
	docker-compose down -v
	@echo "Cleaned up volumes and containers"

# ─── Build Services ──────────────────────────────────────────
build-services:
	cd services/user-service && go build -o bin/user-service ./cmd/main.go
	cd services/group-service && go build -o bin/group-service ./cmd/main.go
	cd services/expense-service && go build -o bin/expense-service ./cmd/main.go
	cd services/settlement-service && go build -o bin/settlement-service ./cmd/main.go
