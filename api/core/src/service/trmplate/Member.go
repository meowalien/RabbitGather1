package member

import "context"

// Member will handle everything about member control
type Member struct {
	*HTTP
}

func (m *Member) Initialize(ctx context.Context) error {
	m.HTTP = &HTTP{}
	return nil
}
