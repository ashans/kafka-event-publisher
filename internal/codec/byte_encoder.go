package codec

import "context"

type ByteEncoder struct {
}

func (b ByteEncoder) Encode(_ context.Context, _ map[string]string, input []byte) ([]byte, error) {
	return input, nil
}
