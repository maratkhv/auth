# Authentication Microservice (gRPC)

![Go](https://img.shields.io/badge/Go-1.20%2B-blue)
![gRPC](https://img.shields.io/badge/gRPC-Protobuf-orange)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15%2B-blue)
![Redis](https://img.shields.io/badge/Redis-7%2B-red)
![Docker](https://img.shields.io/badge/Docker-Compose-yellowgreen)

## ! Designed as a learning project.
An authentication microservice built with Go and gRPC, with postgres and redis used as a storage.

## Main Features

```protobuf
service Auth {
    rpc Register(RegisterRequest) returns (RegisterResponse);  // User registration
    rpc Login(LoginRequest) returns (LoginResponse);           // JWT token generation
    rpc GetPublicKey(GetPublicKeyRequest) returns (GetPublicKeyResponse); // Public key for token verification
}
