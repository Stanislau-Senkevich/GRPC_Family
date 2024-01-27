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

// DenyInvite denies the invite with the specified invite ID.
// It logs information about the operation, such as attempting to deny the invite and whether the operation was successful.
func (s *serverAPI) DenyInvite(
	ctx context.Context,
	req *famv1.DenyInviteRequest,
) (*famv1.DenyInviteResponse, error) {
	const op = "invite.grpc.DenyInvite"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("trying to deny invite",
		slog.Int64("invite_id", req.GetInviteId()))

	err := s.invite.DenyInvite(ctx, req.InviteId)
	if err != nil {
		log.Error("failed to deny invite", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("invite is successfully denied")

	return &famv1.DenyInviteResponse{
		Succeed: true,
	}, nil
}
