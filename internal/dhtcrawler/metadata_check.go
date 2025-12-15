package dhtcrawler

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

func (c *crawler) shouldSkipMetadataDownload(infoHash protocol.ID) bool {
	if c.bloomFilter == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	exists, err := c.bloomFilter.Exists(ctx, infoHash)
	if err != nil {
		c.logger.Warnw("bloom filter query failed; assuming hash is new", "infoHash", infoHash.String(), "error", err)
		return false
	}
	return exists
}
