package family

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

// LeaveFamily allows a user to leave a family identified by the given family ID.
// It logs information about the operation, such as leaving the family and removing the family from the user's family list.
func (s *serverAPI) LeaveFamily(
	ctx context.Context,
	req *famv1.LeaveFamilyRequest,
) (*famv1.LeaveFamilyResponse, error) {
	const op = "family.grpc.LeaveFamily"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("leaving family",
		slog.Int64("family_id", req.GetFamilyId()))

	userID, err := s.family.LeaveFamily(ctx, req.GetFamilyId())
	if errors.Is(err, grpcerror.ErrFamilyNotFound) {
		log.Warn("failed to find family")
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrFamilyNotFound.Error())
	}
	if errors.Is(err, grpcerror.ErrUserNotFound) {
		log.Warn("failed to find user in family")
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrUserNotFound.Error())
	}
	if err != nil {
		log.Error("failed to leave family", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("successfully leaved the family")

	err = s.sso.RemoveFamilyFromList(userID, req.GetFamilyId())
	if err != nil {
		log.Error("failed to delete family from user's family list", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("successfully deleted family from user's family list")

	return &famv1.LeaveFamilyResponse{
		Succeed: true,
	}, nil
}
