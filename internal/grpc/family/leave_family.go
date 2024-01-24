package family

import (
	"context"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
)

func (s *serverAPI) LeaveFamily(
	ctx context.Context,
	req *famv1.LeaveFamilyRequest,
) (*famv1.LeaveFamilyResponse, error) {
	panic("impl me")
}
