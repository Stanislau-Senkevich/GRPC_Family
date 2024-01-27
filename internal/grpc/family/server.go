package family

import (
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/services"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"google.golang.org/grpc"
	"log/slog"
)

type serverAPI struct {
	famv1.UnimplementedFamilyServer
	log    *slog.Logger
	family services.Family
	sso    services.SSO
}

// Register associates the gRPC implementation of the Auth service with the provided gRPC server.
func Register(gRPC *grpc.Server, log *slog.Logger, family services.Family, sso services.SSO) {
	famv1.RegisterFamilyServer(gRPC, &serverAPI{
		log:    log,
		family: family,
		sso:    sso,
	})
}
