package dhtcrawler

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"go.uber.org/zap"
)

const (
	tfileMagic = "TorrentBlobV2"
)

var (
	torrentPrefix = []byte("d4:info")
	torrentSuffix = []byte("e")
)

type tfileWriter struct {
	mu sync.Mutex

	dir string

	logger *zap.SugaredLogger

	maxTorrentsPerFile int
	tfilePrefix        string

	tmpPath   string
	finalPath string
	f         *os.File
	w         *bufio.Writer

	writtenInFile int
}

func newTFileWriter(dir string, maxTorrentsPerFile int, tfilePrefix string, logger *zap.SugaredLogger) *tfileWriter {
	if maxTorrentsPerFile <= 0 {
		maxTorrentsPerFile = 10_000
	}
	return &tfileWriter{
		dir:                dir,
		maxTorrentsPerFile: maxTorrentsPerFile,
		tfilePrefix:        tfilePrefix,
		logger:             logger,
	}
}

func (w *tfileWriter) WriteTorrent(infoHash protocol.ID, rawMetaInfo []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.f == nil {
		if err := w.openNewFileLocked(); err != nil {
			return err
		}
	}

	torrentLen := len(torrentPrefix) + len(rawMetaInfo) + len(torrentSuffix)
	if torrentLen > int(^uint32(0)) {
		return fmt.Errorf("torrent too large: %d bytes", torrentLen)
	}

	var header [len(tfileMagic) + 4 + 20]byte
	copy(header[:], tfileMagic)
	binary.LittleEndian.PutUint32(header[len(tfileMagic):], uint32(torrentLen))
	copy(header[len(tfileMagic)+4:], infoHash[:])

	if _, err := w.w.Write(header[:]); err != nil {
		return err
	}
	if _, err := w.w.Write(torrentPrefix); err != nil {
		return err
	}
	if _, err := w.w.Write(rawMetaInfo); err != nil {
		return err
	}
	if _, err := w.w.Write(torrentSuffix); err != nil {
		return err
	}

	w.writtenInFile++
	if w.writtenInFile >= w.maxTorrentsPerFile {
		if err := w.closeAndFinalizeLocked(); err != nil {
			return err
		}
	}

	return nil
}

func (w *tfileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.closeAndFinalizeLocked()
}

func (w *tfileWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.f == nil {
		return nil
	}
	return w.w.Flush()
}

func (w *tfileWriter) closeAndFinalizeLocked() error {
	if w.f == nil {
		return nil
	}

	var writeErr error
	if err := w.w.Flush(); err != nil {
		writeErr = errors.Join(writeErr, err)
	}
	if err := w.f.Sync(); err != nil {
		writeErr = errors.Join(writeErr, err)
	}
	if err := w.f.Close(); err != nil {
		writeErr = errors.Join(writeErr, err)
	}
	w.f = nil
	w.w = nil
	w.writtenInFile = 0

	if writeErr != nil {
		_ = os.Remove(w.tmpPath)
		return writeErr
	}

	if err := moveFileNoOverwrite(w.tmpPath, w.finalPath); err != nil {
		_ = os.Remove(w.tmpPath)
		return err
	}

	return nil
}

func (w *tfileWriter) openNewFileLocked() error {
	prefix := sanitizeFilenameComponent(w.tfilePrefix)
	if prefix == "unknown" && strings.TrimSpace(w.tfilePrefix) == "" {
		host, err := os.Hostname()
		if err != nil {
			host = "unknown-host"
		}
		prefix = sanitizeFilenameComponent(host)
	}

	stamp := time.Now().Format("20060102-1504")
	base := fmt.Sprintf("%s-%s", prefix, stamp)

	dir := strings.TrimSpace(w.dir)
	if dir == "" {
		dir = "."
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	for attempt := 0; attempt < 1000; attempt++ {
		name := base
		if attempt > 0 {
			name = fmt.Sprintf("%s-%03d", base, attempt)
		}
		final := name + ".tfile"
		tmp := final + ".tmp"

		finalPath := filepath.Join(dir, final)
		tmpPath := filepath.Join(dir, tmp)

		f, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				continue
			}
			return err
		}
		w.f = f
		w.w = bufio.NewWriterSize(f, 1<<20)
		w.tmpPath = filepath.Clean(tmpPath)
		w.finalPath = filepath.Clean(finalPath)
		w.writtenInFile = 0
		if w.logger != nil {
			w.logger.Infow("opened tfile", "finalPath", w.finalPath)
		}
		return nil
	}
	return fmt.Errorf("failed to create unique .tfile.tmp for %q", base)
}

func sanitizeFilenameComponent(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "unknown"
	}
	s = strings.ReplaceAll(s, string(os.PathSeparator), "_")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	s = strings.ReplaceAll(s, ":", "_")
	return s
}

func moveFileNoOverwrite(tmpPath, finalPath string) error {
	if err := os.Link(tmpPath, finalPath); err == nil {
		return os.Remove(tmpPath)
	} else if errors.Is(err, os.ErrExist) {
		return fmt.Errorf("refusing to overwrite existing file %q", finalPath)
	}

	if _, statErr := os.Lstat(finalPath); statErr == nil {
		return fmt.Errorf("refusing to overwrite existing file %q", finalPath)
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return statErr
	}
	return os.Rename(tmpPath, finalPath)
}
