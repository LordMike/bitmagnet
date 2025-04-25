package dhtcrawler

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

// runInfoHashTriage receives discovered hashes on the infoHashTriage channel, determines if they should be crawled,
// and forwards them to the appropriate channel. Possible outcomes are:
// 1. The hash is not in the database, so it is forwarded to the getPeers channel to attempt retrieval of the meta info.
// 2. The hash is in the database, but we don't have the full details of the torrent (for example it was imported outside the DHT crawler,
// and so we don't have the files info), so it is forwarded to the getPeers channel to attempt retrieval of the meta info.
// 3. The hash is in the database, but the seeders/leechers are not known or are outdated, so it is forwarded to the scrape channel.
// 4. The hash is in the database and the seeders/leechers are known and up to date, so it is discarded.
func (c *crawler) runInfoHashTriage(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case reqs := <-c.infoHashTriage.Out():
			allHashes := make([]protocol.ID, 0, len(reqs))
			reqMap := make(map[protocol.ID]nodeHasPeersForHash, len(reqs))
			for _, r := range reqs {
				if _, ok := reqMap[r.infoHash]; ok {
					continue
				}
				allHashes = append(allHashes, r.infoHash)
				reqMap[r.infoHash] = r
			}

			for _, infoHash := range allHashes {
				infoHashStr := strings.ToUpper(infoHash.String())

				// Use the first byte of the infoHash as a directory prefix
				filePath := filepath.Join(c.saveTorrentsRoot, infoHashStr[:2], infoHashStr+".torrent")

				// Check if the file exists
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					// File does not exist, forward to getPeers
					select {
					case <-ctx.Done():
						return
					case c.getPeers.In() <- reqMap[infoHash]:
						continue
					}
				}
			}
		}
	}
}

type triageResult struct {
	InfoHash    protocol.ID
	FilesStatus model.FilesStatus
	FilesCount  model.NullUint
	Seeders     model.NullUint
	Leechers    model.NullUint
	UpdatedAt   time.Time
}
