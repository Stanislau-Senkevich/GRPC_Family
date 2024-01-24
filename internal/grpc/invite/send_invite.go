package invite

import (
	"context"
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
)

func (s *serverAPI) SendInvite(
	ctx context.Context,
	req *famv1.SendInviteRequest,
) (*famv1.SendInviteResponse, error) {
	panic("impl me")
}
