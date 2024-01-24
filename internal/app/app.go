package app

import (
	"fmt"
	grpcapp "github.com/Stanislau-Senkevich/GRPC_Family/internal/app/grpc"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/config"
	jwtmanager "github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/repository/mongodb"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/services/family"
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
	repo, err := mongodb.InitMongoRepository(&cfg.Mongo, log)
	if err != nil {
		panic(fmt.Errorf("failed to initialize repository: %w", err))
	}

	jwtManager := jwtmanager.New(cfg.SigningKey)

	familyService := family.New(log, repo, jwtManager)

	accessibleRoles := map[string][]string{
		"/family.Family/CreateFamily": {"user", "admin"},
	}

	grpcApp := grpcapp.New(
		log, &cfg.GRPC, familyService,
		accessibleRoles, jwtManager,
	)

	return &App{
		GRPCApp: grpcApp,
	}
}
