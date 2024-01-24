package invite

import (
	"context"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
)

func (s *serverAPI) DenyInvite(
	ctx context.Context,
	req *famv1.DenyInviteRequest,
) (*famv1.DenyInviteResponse, error) {
	panic("impl me")
}
