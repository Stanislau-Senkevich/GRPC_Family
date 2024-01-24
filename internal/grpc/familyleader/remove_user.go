package familyleader

import (
	famv1 "github.com/Stanislau-Senkevich/protocols/gen/go/family"
	"golang.org/x/net/context"
)

func (s *serverAPI) RemoveUser(
	ctx context.Context,
	req *famv1.RemoveUserRequest,
) (*famv1.RemoveUserResponse, error) {
	panic("impl me")
}
