package invite

import (
	"context"
	"errors"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/sl"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// SendInvite sends an invitation to a user to join a family.
// It logs information about the operation, such as sending the invite and whether the operation was successful.
func (s *serverAPI) SendInvite(
	ctx context.Context,
	req *famv1.SendInviteRequest,
) (*famv1.SendInviteResponse, error) {
	const op = "invite.grpc.SendInvite"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("sending invite to user",
		slog.Int64("user_id", req.GetUserId()),
		slog.Int64("family_id", req.GetFamilyId()))

	_, err := s.sso.GetUserInfo(req.GetUserId())
	if err != nil {
		log.Warn(grpcerror.ErrUserNotFound.Error())
		return nil, status.Error(codes.Internal, grpcerror.ErrUserNotFound.Error())
	}

	log.Info("invited user exists, trying to send the invite")

	inviteID, err := s.invite.SendInvite(ctx, req.GetFamilyId(), req.GetUserId())
	if errors.Is(err, grpcerror.ErrInviteExist) {
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrInviteExist.Error())
	}
	if errors.Is(err, grpcerror.ErrUserInFamily) {
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserInFamily.Error())
	}
	if errors.Is(err, grpcerror.ErrForbidden) {
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrForbidden.Error())
	}
	if err != nil {
		log.Error("failed to send invite to user", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("invite successfully sent")

	return &famv1.SendInviteResponse{
		InviteId: inviteID,
	}, nil
}
