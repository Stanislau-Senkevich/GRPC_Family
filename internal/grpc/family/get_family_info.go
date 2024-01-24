package family

import (
	"context"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
)

func (s *serverAPI) GetFamilyInfo(
	ctx context.Context,
	req *famv1.GetFamilyInfoRequest,
) (*famv1.GetFamilyInfoResponse, error) {
	panic("impl me")
}
