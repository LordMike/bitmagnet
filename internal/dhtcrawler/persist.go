package dhtcrawler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// runPersistTorrents serially drains the persist queue and writes torrents to disk as they arrive.
func (c *crawler) runPersistTorrents(ctx context.Context) {
	handler := func(item infoHashWithMetaInfo) {
		if !c.saveTorrents {
			return
		}
		if err := c.saveRawMetadataToFile(item.infoHash.String(), item.MetaInfoBytes); err != nil {
			c.logger.Errorw("failed to save torrent", "infoHash", item.infoHash.String(), "error", err)
		}
	}
	if err := c.persistTorrents.Run(ctx, handler); err != nil && !errors.Is(err, context.Canceled) {
		c.logger.Errorw("persist worker stopped unexpectedly", "error", err)
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

	// Create the temporary file exclusively so parallel writers defer to the first
	tempFile, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			c.logger.Debugw("temp file already exists, skipping save", "tempFilePath", tempFilePath)
			return nil
		}
		c.logger.Errorw("failed to create temp file", "tempFilePath", tempFilePath, "error", err)
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempFilePath)
	}()

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
		if _, statErr := os.Stat(finalFilePath); statErr == nil {
			return nil
		}
		c.logger.Errorw("failed to rename temp file to final file", "tempFilePath", tempFilePath, "finalFilePath", finalFilePath, "error", err)
		return fmt.Errorf("failed to rename temp file to final file: %v", err)
	}

	c.logger.Infow("saved torrent", "infoHash", infoHash)
	return nil
}
