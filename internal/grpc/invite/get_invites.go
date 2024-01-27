package invite

import (
	"context"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/sl"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// GetInvites retrieves the invites of the user.
// It logs information about the operation, such as attempting to retrieve the invites and whether the operation was successful.
func (s *serverAPI) GetInvites(
	ctx context.Context,
	_ *famv1.GetInvitesRequest,
) (*famv1.GetInvitesResponse, error) {
	const op = "invite.grpc.GetInvites"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("retrieving invites of user")

	invites, err := s.invite.GetInvites(ctx)
	if err != nil {
		log.Error("failed to retrieve invites", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("invites successfully retrieved")

	return &famv1.GetInvitesResponse{
		Invites: invites,
	}, nil
}
