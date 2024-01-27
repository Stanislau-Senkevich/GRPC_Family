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

// DeleteUserInvites deletes all invites associated with the specified user ID.
// It logs information about the operation, such as attempting to delete user invites and whether the operation was successful.
func (s *serverAPI) DeleteUserInvites(
	ctx context.Context,
	req *famv1.DeleteUserInvitesRequest,
) (*famv1.DeleteUserInvitesResponse, error) {
	const op = "invite.grpc.DeleteUserInvites"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("trying to delete user invites",
		slog.Int64("user_id", req.GetUserId()))

	err := s.invite.DeleteUserInvites(ctx, req.GetUserId())
	if err != nil {
		log.Error("failed to delete user's invites",
			sl.Err(err), slog.Int64("user_id", req.GetUserId()))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("user's invites successfully deleted",
		slog.Int64("user_id", req.GetUserId()))

	return &famv1.DeleteUserInvitesResponse{
		Succeed: true,
	}, nil
}
