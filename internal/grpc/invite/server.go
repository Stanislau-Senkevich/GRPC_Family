package invite

import (
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/services"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"google.golang.org/grpc"
	"log/slog"
)

type serverAPI struct {
	famv1.UnimplementedInviteServer
	log    *slog.Logger
	invite services.Invite
	sso    services.SSO
}

// Register associates the gRPC implementation of the Auth service with the provided gRPC server.
func Register(
	gRPC *grpc.Server,
	log *slog.Logger,
	invite services.Invite,
	sso services.SSO) {
	famv1.RegisterInviteServer(gRPC, &serverAPI{
		log:    log,
		invite: invite,
		sso:    sso,
	})
}
