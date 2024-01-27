package app

import (
	"context"
	"fmt"
	grpcapp "github.com/Stanislau-Senkevich/GRPC_Family/internal/app/grpc"
	grpcclient "github.com/Stanislau-Senkevich/GRPC_Family/internal/client/sso/grpc"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/config"
	jwtmanager "github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/jwt"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/repository/mongodb"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/services/family"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/services/familyleader"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/services/invite"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/services/sso"
	"log/slog"
)

type App struct {
	GRPCAppServer *grpcapp.App
}

// New creates a new instance of the application with the provided configuration and dependencies.
func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	log.Info("starting initialize app")

	repo, err := mongodb.InitMongoRepository(&cfg.Mongo, log)
	if err != nil {
		panic(fmt.Errorf("failed to initialize repository: %w", err))
	}
	log.Info("repository initialized")

	jwtManager := jwtmanager.New(cfg.SigningKey)
	log.Info("jwt-manager initialized")

	ssoClient, err := grpcclient.New(
		context.Background(), log,
		cfg.ClientsConfig.SSO.Address,
		cfg.ClientsConfig.SSO.Timeout,
		cfg.ClientsConfig.SSO.RetriesCount)
	if err != nil {
		panic(fmt.Errorf("failed to initialize client SSO: %w", err))
	}
	log.Info("sso client initialized")

	familyService := family.New(log, repo, jwtManager)
	log.Info("family service initialized")

	leaderService := familyleader.New(log, repo, repo, jwtManager)
	log.Info("family leader service initialized")

	inviteService := invite.New(log, repo, repo, jwtManager)
	log.Info("invite service initialized")

	ssoService := sso.New(ssoClient, jwtManager, cfg.ClientsConfig.AdminEmail, cfg.ClientsConfig.AdminPassword)
	log.Info("sso service initialized")

	accessibleRoles := map[string][]string{
		"/family.Family/CreateFamily":       {"user", "admin"},
		"/family.Family/LeaveFamily":        {"user", "admin"},
		"/family.Family/GetFamilyInfo":      {"user", "admin"},
		"/family.Invite/GetInvites":         {"user", "admin"},
		"/family.Invite/SendInvite":         {"user", "admin"},
		"/family.Invite/AcceptInvite":       {"user", "admin"},
		"/family.Invite/DenyInvite":         {"user", "admin"},
		"family.Invite/DeleteUserInvites":   {"admin"},
		"/family.FamilyLeader/RemoveUser":   {"user", "admin"},
		"/family.FamilyLeader/DeleteFamily": {"user", "admin"},
	}

	grpcApp := grpcapp.New(
		log, &cfg.GRPC,
		familyService, leaderService,
		inviteService, ssoService,
		accessibleRoles, jwtManager,
	)

	log.Info("grpc-server initialized")

	return &App{
		GRPCAppServer: grpcApp,
	}
}
