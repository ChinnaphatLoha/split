FROM golang:1.24-alpine AS builder

RUN apk add --no-cache protobuf protobuf-dev git
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app

COPY services/expense-service/ .
COPY proto/expense/expense.proto /proto/expense.proto

RUN rm -f proto/expense.go proto/expense_grpc.go

RUN protoc --proto_path=/proto --proto_path=/usr/include \
    --go_out=proto --go_opt=paths=source_relative \
    --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
    /proto/expense.proto

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /expense-service ./cmd/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /expense-service /expense-service
EXPOSE 50053
CMD ["/expense-service"]
