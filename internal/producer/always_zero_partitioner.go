package producer

import "context"

type AlwaysZeroPartitioner struct {
}

func (p AlwaysZeroPartitioner) Partition(_ context.Context, _ map[string]string, _, _, _, _ []byte,
	_ map[string]string, _ int32) (int32, error) {
	return 0, nil
}
