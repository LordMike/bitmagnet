package dhtcrawler

import (
	"context"
	"strings"

	"fmt"
	"os"
	"path/filepath"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
)

// runPersistTorrents waits on the persistTorrents channel, and persists torrents to the database in batches.
// After persisting each batch it will publish a message to the classifier,
// and forward the hash on the scrape channel to attempt finding the seeders/leechers.
func (c *crawler) runPersistTorrents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case is := <-c.persistTorrents.Out():
			hashMap := make(map[protocol.ID]infoHashWithMetaInfo, len(is))
			for _, i := range is {
				if _, ok := hashMap[i.infoHash]; ok {
					continue
				}
				hashMap[i.infoHash] = i
			}

			// Persist to disk
			if c.saveTorrents {
				for _, i := range is {
					c.saveRawMetadataToFile(i.infoHash.String(), i.MetaInfoBytes)
				}
			}

			for _, i := range hashMap {
				select {
				case <-ctx.Done():
					return
				case c.scrape.In() <- i.nodeHasPeersForHash:
					continue
				}
			}
		}
	}
}

func (c *crawler) saveRawMetadataToFile(infoHash string, rawMetaInfo []byte) error {
	// Convert infoHash to uppercase to ensure consistency
	infoHash = strings.ToUpper(infoHash)

	// Create a two-level trie directory structure using the first 4 characters of the infoHash
	dir1 := infoHash[:2] // First 2 characters
	directory := filepath.Join(c.saveTorrentsRoot, dir1)

	// Create the directory structure if it doesn't exist
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		c.logger.Errorw("failed to create directory", "directory", directory, "error", err)
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Define the final file path and temporary file path
	finalFilePath := filepath.Join(directory, infoHash+".torrent")
	tempFilePath := finalFilePath + c.saveTorrentsTempSuffix

	// Check if the final file already exists, and skip if it does
	if _, err := os.Stat(finalFilePath); err == nil {
		c.logger.Debugw("File already exists, skipping save", "filePath", finalFilePath)
		return nil
	}

	// Create and write to the temporary file
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		c.logger.Errorw("failed to create temp file", "tempFilePath", tempFilePath, "error", err)
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()

	var writeErr error
	_, writeErr = tempFile.Write([]byte("d4:info"))

	_, err = tempFile.Write(rawMetaInfo)
	if writeErr == nil {
		writeErr = err
	}

	_, err = tempFile.Write([]byte("e"))
	if writeErr == nil {
		writeErr = err
	}

	if writeErr != nil {
		c.logger.Errorw("failed to write raw metadata to temp file", "tempFilePath", tempFilePath, "error", err)
		return fmt.Errorf("failed to write raw metadata to temp file: %v", err)
	}

	// Ensure all data is written to disk
	if err := tempFile.Sync(); err != nil {
		c.logger.Errorw("failed to sync temp file", "tempFilePath", tempFilePath, "error", err)
		return fmt.Errorf("failed to sync temp file: %v", err)
	}

	// Rename the temp file to the final file
	if err := os.Rename(tempFilePath, finalFilePath); err != nil {
		c.logger.Errorw("failed to rename temp file to final file", "tempFilePath", tempFilePath, "finalFilePath", finalFilePath, "error", err)
		return fmt.Errorf("failed to rename temp file to final file: %v", err)
	}

	c.logger.Infof("saved torrent", infoHash)
	return nil
}

func createTorrentModel(
	hash protocol.ID,
	info metainfo.Info,
	savePieces bool,
	saveFilesThreshold uint,
) (model.Torrent, error) {
	name := info.BestName()
	private := false
	if info.Private != nil {
		private = *info.Private
	}
	var filesCount model.NullUint
	filesStatus := model.FilesStatusSingle
	if len(info.Files) > 0 {
		filesStatus = model.FilesStatusMulti
		filesCount = model.NewNullUint(uint(len(info.Files)))
	}
	var files []model.TorrentFile
	for i, file := range info.Files {
		if i >= int(saveFilesThreshold) {
			filesStatus = model.FilesStatusOverThreshold
			break
		}
		files = append(files, model.TorrentFile{
			InfoHash: hash,
			Index:    uint(i),
			Path:     file.DisplayPath(&info),
			Size:     uint(file.Length),
		})
	}
	var pieces model.TorrentPieces
	if savePieces {
		pieces = model.TorrentPieces{
			InfoHash:    hash,
			PieceLength: info.PieceLength,
			Pieces:      info.Pieces,
		}
	}
	return model.Torrent{
		InfoHash:    hash,
		Name:        name,
		Size:        uint(info.TotalLength()),
		Private:     private,
		Pieces:      pieces,
		Files:       files,
		FilesStatus: filesStatus,
		FilesCount:  filesCount,
		Sources: []model.TorrentsTorrentSource{
			{
				Source:   "dht",
				InfoHash: hash,
			},
		},
	}, nil
}

const classifyBatchSize = 100

// runPersistSources waits on the persistSources channel for scraped torrents, and persists sources
// (which includes discovery date, seeders and leechers) to the database in batches.
func (c *crawler) runPersistSources(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		// case scrapes := <-c.persistSources.Out():
		case <-c.persistSources.Out():
			// srcs := make([]*model.TorrentsTorrentSource, 0, len(scrapes))
			// hashSet := make(map[protocol.ID]struct{}, len(scrapes))
			// for _, s := range scrapes {
			// 	if _, ok := hashSet[s.infoHash]; ok {
			// 		continue
			// 	}
			// 	hashSet[s.infoHash] = struct{}{}
			// 	if src, err := createTorrentSourceModel(s); err != nil {
			// 		c.logger.Errorf("error creating torrent source model: %s", err.Error())
			// 	} else {
			// 		srcs = append(srcs, &src)
			// 	}
			// }
			// if persistErr := c.dao.WithContext(ctx).TorrentsTorrentSource.Clauses(
			// 	clause.OnConflict{
			// 		Columns: []clause.Column{
			// 			{Name: string(c.dao.TorrentsTorrentSource.InfoHash.ColumnName())},
			// 			{Name: string(c.dao.TorrentsTorrentSource.Source.ColumnName())},
			// 		},
			// 		DoUpdates: clause.AssignmentColumns([]string{
			// 			string(c.dao.TorrentsTorrentSource.Seeders.ColumnName()),
			// 			string(c.dao.TorrentsTorrentSource.Leechers.ColumnName()),
			// 			// sets to null, fixes torrents indexed before 0.8.0 with published_at 0001-01-01 00:00:00+00:
			// 			string(c.dao.TorrentsTorrentSource.PublishedAt.ColumnName()),
			// 			string(c.dao.TorrentsTorrentSource.UpdatedAt.ColumnName()),
			// 		}),
			// 	},
			// ).Where(
			// 	// check that the torrent record hasn't been deleted:
			// 	gen.Exists(c.dao.WithContext(ctx).Torrent.Where(
			// 		c.dao.Torrent.InfoHash.EqCol(c.dao.TorrentsTorrentSource.InfoHash),
			// 	)),
			// ).CreateInBatches(srcs, 100); persistErr != nil {
			// 	c.logger.Errorf("error persisting torrent sources: %s", persistErr.Error())
			// } else {
			// 	c.persistedTotal.With(prometheus.Labels{"entity": "TorrentsTorrentSource"}).Add(float64(len(srcs)))
			// 	c.logger.Debugw("persisted torrent sources", "count", len(srcs))
			// }
			return
		}
	}
}

func createTorrentSourceModel(
	result infoHashWithScrape,
) (model.TorrentsTorrentSource, error) {
	seeders := model.NewNullUint(uint(result.bfsd.ApproximatedSize()))
	leechers := model.NewNullUint(uint(result.bfpe.ApproximatedSize()))
	return model.TorrentsTorrentSource{
		Source:   "dht",
		InfoHash: result.infoHash,
		Seeders:  seeders,
		Leechers: leechers,
	}, nil
}
