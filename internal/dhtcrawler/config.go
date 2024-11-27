package dhtcrawler

import (
	"time"

	adht "github.com/anacrolix/dht/v2"
)

type Config struct {
	// ScalingFactor is a rough proxy for resource usage of the crawler; concurrency and buffer size of the various
	// pipeline channels are multiplied by this value. Diminishing returns may result from exceeding the default value of 10.
	// Since the software has not been tested on a wide variety of hardware and network conditions your mileage may vary here...
	ScalingFactor                uint
	BootstrapNodes               []string
	ReseedBootstrapNodesInterval time.Duration
	// SaveFilesThreshold specifies a maximum number of files in a torrent before file information is discarded.
	// Some torrents contain thousands of files which can severely impact performance and uses a lot of disk space.
	SaveFilesThreshold uint
	// SavePieces when true, torrent pieces will be persisted to the database.
	// The pieces take up quite a lot of space, and aren't currently very useful, but they may be used by future features.
	SavePieces bool
	// RescrapeThreshold is the amount of time that must pass before a torrent is rescraped to count seeders and leechers.
	RescrapeThreshold time.Duration

	SaveTorrents     bool
	SaveTorrentsRoot string
	// When multiple instances of dht_crawler are running, it is possible to avoid torrent corruption by setting a unique temp file suffix for each instance
	SaveTorrentsTempSuffix string
}

func NewDefaultConfig() Config {
	return Config{
		ScalingFactor:                10,
		BootstrapNodes:               defaultBootstrapNodes,
		ReseedBootstrapNodesInterval: time.Minute,
		SaveFilesThreshold:           100,
		SavePieces:                   false,
		RescrapeThreshold:            time.Hour * 24 * 30,
		SaveTorrents:                 false,
		SaveTorrentsRoot:             "./torrents",
		SaveTorrentsTempSuffix:       ".tmp",
	}
}

// https://github.com/anacrolix/dht/blob/92b36a3fa7a37a15e08684337b47d8d0fb322ab6/dht.go#L106
var defaultBootstrapNodes = []string{
	"router.utorrent.com:6881",
	"router.bittorrent.com:6881",
	"dht.transmissionbt.com:6881",
	"dht.aelitis.com:6881",     // Vuze
	"router.silotis.us:6881",   // IPv6
	"dht.libtorrent.org:25401", // @arvidn's
}
