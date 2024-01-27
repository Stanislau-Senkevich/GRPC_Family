package family

import (
	"context"
	"errors"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/domain/models"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/sl"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// GetFamilyInfo retrieves information about the members of a family identified by the given family ID.
// It logs information about the operation, such as retrieving family members' information and handling any errors.
func (s *serverAPI) GetFamilyInfo(
	ctx context.Context,
	req *famv1.GetFamilyInfoRequest,
) (*famv1.GetFamilyInfoResponse, error) {
	const op = "family.grpc.GetFamilyInfo"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("getting family info",
		slog.Int64("family_id", req.GetFamilyId()))

	IDs, err := s.family.GetFamilyMembersIDs(ctx, req.GetFamilyId())
	if errors.Is(err, grpcerror.ErrFamilyNotFound) {
		log.Warn(grpcerror.ErrFamilyNotFound.Error(),
			slog.Int64("family_id", req.GetFamilyId()))
		return nil, status.Error(codes.InvalidArgument, grpcerror.ErrFamilyNotFound.Error())
	}
	if errors.Is(err, grpcerror.ErrForbidden) {
		return nil, status.Error(codes.PermissionDenied, grpcerror.ErrForbidden.Error())
	}
	if err != nil {
		log.Error("failed to get members' id", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("successfully got family members' ids")

	info := make([]*famv1.UserInfo, 0, len(IDs))

	damaged := make([]int64, 0)

	for _, userID := range IDs {
		user, err := s.sso.GetUserInfo(userID)
		if err != nil {
			log.Error("user's info is missing",
				sl.Err(err), slog.Int64("user_id", userID))
			damaged = append(damaged, userID)
			continue
		}

		info = append(info, models.ConvertToInfo(user))
	}

	resp := &famv1.GetFamilyInfoResponse{
		Info: info,
	}

	if len(damaged) > 0 {
		return resp, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	return resp, nil
}
