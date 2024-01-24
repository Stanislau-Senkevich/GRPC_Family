package app

import (
	"fmt"
	grpcapp "github.com/Stanislau-Senkevich/GRPC_Family/internal/app/grpc"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/config"
	jwtmanager "github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/repository/mongodb"
	"log/slog"
)

type App struct {
	GRPCApp *grpcapp.App
}

// New creates a new instance of the application with the provided configuration and dependencies.
func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	_, err := mongodb.InitMongoRepository(&cfg.Mongo, log)
	if err != nil {
		panic(fmt.Errorf("failed to initialize repository: %w", err))
	}

	jwtManager := jwtmanager.New(cfg.SigningKey)

	accessibleRoles := map[string][]string{}

	grpcApp := grpcapp.New(
		log, &cfg.GRPC,
		accessibleRoles, jwtManager,
	)

	return &App{
		GRPCApp: grpcApp,
	}
}
