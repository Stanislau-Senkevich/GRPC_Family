# GRPC Family Microservice

---------------

## Introduction

This API based on gRPC technology provides interface for creating and administrating user's unions(families)
with integration with my another GRPC microservice (https://github.com/Stanislau-Senkevich/GRPC_SSO).

To test this API on your own you should download protocols from https://github.com/Stanislau-Senkevich/protocols
and send grpc requests on <span style="color: blue"> grpc://droplet.senkevichdev.work:33033 </span> 

### Models
- Admin
- User
- Family
- Invite


#### User
- Can create families and become its leader.
- Leader of family is allowed to send invitations to family to another users. He also allowed to kick users from families or delete a whole family.
- Other members of family can check info about users in family and can leave family, if necessary
- Users also can accept or deny invitations to other families which were sent to them.

#### Admin
- All user's features
- Allowed to operate with families same as its leaders

------------------
## Technologies
- #### Go 1.21
- #### gRPC
- #### MongoDB
- #### Docker
- #### Kubernetes
- #### JWT-tokens
- #### DNS
- #### CI/CD (GitHub Actions)

-----------------
## Realization features
- #### Microservice architecture 
- #### Clean architecture
- #### Functional tests for handlers
- #### Linter
- #### Logging with slog package

-----------------
### Tools and libraries

### Protocols and Utilities

- `Stanislau-Senkevich/protocols`: Protocol implementations.

### Authentication

- `golang-jwt/jwt`: JWT functionality.

### Middleware

- `grpc-ecosystem/go-grpc-middleware/v2`: provides gRPC middleware.

### Configuration

- `ilyakaznacheev/cleanenv`: Environment variable reading and validation.
- `spf13/viper`: Configuration solution for Go applications.
- `subosito/gotenv`: Loading environment variables from .env files.

### Database

- `go.mongodb.org/mongo-driver`: Go package providing driver and functinality to interact with MongoDB.

### gRPC

- `google.golang.org/grpc`: gRPC Go implementation.

### Logging

- `log/slog`: standard Go library for logging.

### Protocol Buffers

- `google.golang.org/protobuf`: Protocol Buffers serialization.

