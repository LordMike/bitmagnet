package dhtcrawler

import (
	"context"
	"errors"
)

// runPersistTorrents serially drains the persist queue and writes torrents to disk as they arrive.
func (c *crawler) runPersistTorrents(ctx context.Context) {
	defer func() {
		if c.tfileWriter != nil {
			_ = c.tfileWriter.Close()
		}
		if c.persistDone != nil {
			close(c.persistDone)
		}
	}()
	handler := func(item infoHashWithMetaInfo) {
		if err := c.tfileWriter.WriteRawMetadata(item.MetaInfoBytes); err != nil {
			c.logger.Errorw("failed to persist torrent", "infoHash", item.infoHash.String(), "error", err)
		}
	}
	if err := c.persistTorrents.Run(ctx, handler); err != nil && !errors.Is(err, context.Canceled) {
		c.logger.Errorw("persist worker stopped unexpectedly", "error", err)
	}
}
