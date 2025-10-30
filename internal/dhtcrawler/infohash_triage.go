package dhtcrawler

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

// runInfoHashTriage receives discovered hashes, deduplicates them, and forwards only those without a saved .torrent file
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

				infoHashStr := strings.ToUpper(r.infoHash.String())
				dir1 := infoHashStr[:2]
				finalFilePath := filepath.Join(c.saveTorrentsRoot, dir1, infoHashStr+".torrent")

				if _, err := os.Stat(finalFilePath); err == nil {
					continue
				} else if !os.IsNotExist(err) {
					c.logger.Warnw("skipping infohash due to stat error", "infoHash", infoHashStr, "error", err)
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
