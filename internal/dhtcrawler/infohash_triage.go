package dhtcrawler

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

// runInfoHashTriage receives discovered hashes, deduplicates them, and forwards those which require metadata download
// to the getPeers stage so that metadata can be fetched.
func (c *crawler) runInfoHashTriage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case reqs := <-c.infoHashTriage.Out():
			seen := make(map[protocol.ID]struct{}, len(reqs))
			for _, r := range reqs {
				if _, ok := seen[r.infoHash]; ok {
					continue
				}
				seen[r.infoHash] = struct{}{}

				infoHashStr := r.infoHash.String()
				if c.shouldSkipMetadataDownload(r.infoHash) {
					continue
				}

				select {
				case <-ctx.Done():
					return
				case c.getPeers.In() <- r:
					c.logger.Debugw("queued get_peers", "infoHash", infoHashStr)
				}
			}
		}
	}
}
