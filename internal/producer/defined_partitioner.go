package producer

import (
	"context"
	"ev_pub/internal/errors"
	"strconv"
)

type DefinedPartitioner struct {
}

func (p DefinedPartitioner) Partition(_ context.Context, options map[string]string, _, _, _, _ []byte,
	_ map[string]string, partitionCount int32) (int32, error) {
	partition, err := strconv.ParseInt(options["partition"], 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, `cannot parse partition option "partition"`)
	}
	if partition < 0 {
		return 0, errors.New("invalid partition number: " + strconv.FormatInt(partition, 10))
	}
	if partition >= int64(partitionCount) {
		return 0, errors.New("partition number should be <: " + strconv.FormatInt(int64(partitionCount), 10))
	}
	return int32(partition), nil
}
