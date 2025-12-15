package dhtcrawler

import "github.com/bitmagnet-io/bitmagnet/internal/protocol"

func (c *crawler) shouldSkipMetadataDownload(infoHash protocol.ID) bool {
	return false
}
