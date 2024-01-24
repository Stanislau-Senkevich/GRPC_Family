package familyleader

import (
	"context"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
)

func (s *serverAPI) DeleteFamily(
	ctx context.Context,
	req *famv1.DeleteFamilyRequest,
) (*famv1.DeleteFamilyResponse, error) {
	panic("impl me")
}
