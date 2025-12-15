package dhtcrawler

import (
	"time"
)

type Config struct {
	// ScalingFactor is a rough proxy for resource usage of the crawler; concurrency and buffer size of the various
	// pipeline channels are multiplied by this value. Diminishing returns may result from exceeding the default value of 10.
	// Since the software has not been tested on a wide variety of hardware and network conditions your mileage may vary here...
	ScalingFactor                uint
	BootstrapNodes               []string
	ReseedBootstrapNodesInterval time.Duration
	SaveTorrentsRoot             string
}

func NewDefaultConfig() Config {
	return Config{
		ScalingFactor:                10,
		BootstrapNodes:               defaultBootstrapNodes,
		ReseedBootstrapNodesInterval: time.Minute,
		SaveTorrentsRoot:             "./torrents",
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
