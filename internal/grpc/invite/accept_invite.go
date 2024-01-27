package invite

import (
	"context"
	"errors"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// AcceptInvite accepts an invitation with the given invite ID and adds the associated family to the user's family list.
// It logs information about the operation, such as attempting to accept the invite and adding the family to the user's family list.
func (s *serverAPI) AcceptInvite(
	ctx context.Context,
	req *famv1.AcceptInviteRequest,
) (*famv1.AcceptInviteResponse, error) {
	const op = "invite.grpc.AcceptInvite"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("trying to accept invite",
		slog.Int64("invite_id", req.GetInviteId()))

	familyID, err := s.invite.AcceptInvite(ctx, req.GetInviteId())
	if errors.Is(err, grpcerror.ErrInviteNotFound) {
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrInviteNotFound.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("invite accepted, trying to add family to user's family list")

	err = s.sso.AddFamilyToList(ctx, familyID)
	if err != nil {
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("family added to user's family list")

	return &famv1.AcceptInviteResponse{
		FamilyId: familyID,
	}, nil
}
