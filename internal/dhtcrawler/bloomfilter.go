package dhtcrawler

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type bloomFilter interface {
	Exists(ctx context.Context, infoHash protocol.ID) (bool, error)
	Add(ctx context.Context, infoHash protocol.ID) error
	Close() error
}
