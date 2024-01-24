package invite

import (
	"context"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
)

func (s *serverAPI) AcceptInvite(
	ctx context.Context,
	req *famv1.AcceptInviteRequest,
) (*famv1.AcceptInviteResponse, error) {
	panic("impl me")
}
