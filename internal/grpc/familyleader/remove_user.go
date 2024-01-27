package familyleader

import (
	"errors"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/sl"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// RemoveUser removes a user with the given user ID from the specified family.
// It logs information about the operation, such as attempting to remove the user from the family and removing the family from the user's family list.
func (s *serverAPI) RemoveUser(
	ctx context.Context,
	req *famv1.RemoveUserRequest,
) (*famv1.RemoveUserResponse, error) {
	const op = "family.grpc.RemoveUser"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("trying to remove user from family",
		slog.Int64("family_id", req.GetFamilyId()),
		slog.Int64("user_id", req.GetUserId()))

	err := s.familyLeader.RemoveUserFromFamily(ctx,
		req.GetFamilyId(), req.GetUserId())
	if errors.Is(err, grpcerror.ErrUserNotFound) {
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotFound.Error())
	}
	if errors.Is(err, grpcerror.ErrUserNotInFamily) {
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotInFamily.Error())
	}
	if errors.Is(err, grpcerror.ErrFamilyNotFound) {
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrFamilyNotFound.Error())
	}
	if errors.Is(err, grpcerror.ErrForbidden) {
		return nil, status.Error(codes.PermissionDenied, grpcerror.ErrForbidden.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("user removed from family")

	log.Info("trying to remove family from user's list")

	err = s.sso.RemoveFamilyFromList(req.GetUserId(), req.GetFamilyId())
	if err != nil {
		log.Error("failed to remove family from user's list", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	return &famv1.RemoveUserResponse{
		UserId: req.GetUserId(),
	}, nil
}
