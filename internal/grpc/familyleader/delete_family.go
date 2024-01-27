package familyleader

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

// DeleteFamily deletes the family with the given family ID from the system.
// It logs information about the operation, such as attempting to delete the family and removing the family from users' family lists.
func (s *serverAPI) DeleteFamily(
	ctx context.Context,
	req *famv1.DeleteFamilyRequest,
) (*famv1.DeleteFamilyResponse, error) {
	const op = "familyleader.grpc.DeleteFamily"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("trying to delete family",
		slog.Int64("family_id", req.GetFamilyId()))

	members, err := s.familyLeader.DeleteFamily(ctx, req.GetFamilyId())
	if errors.Is(err, grpcerror.ErrFamilyNotFound) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, grpcerror.ErrForbidden) {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("family deleted, removing its id from users' lists")

	for _, id := range members {
		err = s.sso.RemoveFamilyFromList(id, req.GetFamilyId())
		if err != nil {
			log.Warn("failed to delete family from user's family list",
				slog.Int64("user_id", id), sl.Err(err))
		}
	}

	log.Info("family deleted",
		slog.Int64("family_id", req.GetFamilyId()))

	return &famv1.DeleteFamilyResponse{
		Succeed: true,
	}, nil
}
