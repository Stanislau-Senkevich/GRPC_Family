package family

import (
	"context"
	grpcerror "github.com/Stanislau-Senkevich/GRPC_Family/internal/error"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/sl"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

// CreateFamily creates a new family and also adds it to the user's family list.
// It logs information about the operation, such as creating the family and adding it to the user's list.
func (s *serverAPI) CreateFamily(
	ctx context.Context,
	_ *famv1.CreateFamilyRequest,
) (*famv1.CreateFamilyResponse, error) {
	const op = "family.grpc.CreateFamily"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("creating family")

	familyID, err := s.family.CreateFamily(ctx)
	if err != nil {
		log.Error("failed to create family", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("family created", slog.Int64("family_id", familyID))

	if err = s.sso.AddFamilyToList(ctx, familyID); err != nil {
		log.Error("failed to add family to user's family list", sl.Err(err))
		return nil, status.Error(codes.Internal, grpcerror.ErrInternalError.Error())
	}

	log.Info("family added to user's family list")

	return &famv1.CreateFamilyResponse{
		FamilyId: familyID,
	}, nil
}
