package producer

import (
	"context"
	"fmt"
	"strconv"
)

type DefinedPartitioner struct {
}

func (p DefinedPartitioner) Partition(_ context.Context, options map[string]string, _, _, _, _ []byte,
	_ map[string]string, partitionCount int32) (int32, error) {
	partition, err := strconv.ParseInt(options["partition"], 10, 32)
	if err != nil {
		return 0, err
	}
	if partition < 0 {
		return 0, fmt.Errorf("invalid partition number: %d", partition)
	}
	if partition >= int64(partitionCount) {
		return 0, fmt.Errorf("partition number should be <: %d", partitionCount)
	}
	return int32(partition), nil
}
