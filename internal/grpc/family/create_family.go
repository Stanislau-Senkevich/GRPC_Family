package family

import (
	"context"
	"github.com/Stanislau-Senkevich/GRPC_Family/internal/lib/sl"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

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
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Info("family created", slog.Int64("family_id", familyID))

	return &famv1.CreateFamilyResponse{
		FamilyId: familyID,
	}, nil
}
