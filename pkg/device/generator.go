package device

import "context"

type KeyGenerator interface {
	Generate(ctx context.Context) ([]byte, []byte, error)
}
