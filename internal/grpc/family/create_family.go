package family

import (
	"context"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
)

func (s *serverAPI) CreateFamily(
	ctx context.Context,
	req *famv1.CreateFamilyRequest,
) (*famv1.CreateFamilyResponse, error) {
	panic("impl me")
}
