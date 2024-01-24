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

func (s *serverAPI) LeaveFamily(
	ctx context.Context,
	req *famv1.LeaveFamilyRequest,
) (*famv1.LeaveFamilyResponse, error) {
	const op = "family.grpc.CreateFamily"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("leaving family")

	err := s.family.LeaveFamily(ctx, req.FamilyId)
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
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Info("successfully leaved the family")

	return &famv1.LeaveFamilyResponse{
		Succeed: true,
	}, nil
}
