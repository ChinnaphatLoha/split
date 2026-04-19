FROM golang:1.24-alpine AS builder

RUN apk add --no-cache protobuf protobuf-dev git
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app

# Copy service source code
COPY services/user-service/ .

# Copy proto file (flat, not nested in user/ subdir)
COPY proto/user/user.proto /proto/user.proto

# Remove hand-written stubs (will be replaced by protoc output)
RUN rm -f proto/user.go proto/user_grpc.go

# Generate proper protobuf Go code into proto/ directory
RUN protoc --proto_path=/proto --proto_path=/usr/include \
    --go_out=proto --go_opt=paths=source_relative \
    --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
    /proto/user.proto

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./cmd/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /user-service /user-service
EXPOSE 50051
CMD ["/user-service"]
