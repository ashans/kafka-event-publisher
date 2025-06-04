package producer

import (
	"context"
)

type Partitioner interface {
	Partition(ctx context.Context, options map[string]string, inputKey, inputVal, encodedKey, encodedVal []byte,
		headers map[string]string, partitionCount int32) (int32, error)
}
