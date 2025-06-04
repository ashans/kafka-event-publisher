package codec

import (
	"context"
)

type Encoder interface {
	Encode(ctx context.Context, options map[string]string, input []byte) ([]byte, error)
}
