package dhtcrawler

import (
	"context"
	"path/filepath"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
	boom "github.com/tylertreat/BoomFilters"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Config            Config
	KTable            ktable.Table
	Client            lazy.Lazy[client.Client]
	MetainfoRequester metainforequester.Requester
	DiscoveredNodes   concurrency.BatchingChannel[ktable.Node] `name:"dht_discovered_nodes"`
	Logger            *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Worker worker.Worker `group:"workers"`
}

func New(params Params) Result {
	var c crawler
	return Result{
		Worker: worker.NewWorker(
			"dht_crawler",
			fx.Hook{
				OnStart: func(context.Context) error {
					scalingFactor := int(params.Config.ScalingFactor)
					cl, err := params.Client.Get()
					if err != nil {
						return err
					}
					c = crawler{
						kTable:                       params.KTable,
						client:                       cl,
						metainfoRequester:            params.MetainfoRequester,
						bootstrapNodes:               params.Config.BootstrapNodes,
						reseedBootstrapNodesInterval: params.Config.ReseedBootstrapNodesInterval,
						getOldestNodesInterval:       time.Second * 10,
						oldPeerThreshold:             time.Minute * 15,
						discoveredNodes:              params.DiscoveredNodes,
						nodesForPing:                 concurrency.NewBufferedConcurrentChannel[ktable.Node](scalingFactor, scalingFactor),
						nodesForFindNode:             concurrency.NewBufferedConcurrentChannel[ktable.Node](10*scalingFactor, 10*scalingFactor),
						nodesForSampleInfoHashes:     concurrency.NewBufferedConcurrentChannel[ktable.Node](10*scalingFactor, 10*scalingFactor),
						infoHashTriage:               concurrency.NewBatchingChannel[nodeHasPeersForHash](10*scalingFactor, 1000, 20*time.Second),
						getPeers:                     concurrency.NewBufferedConcurrentChannel[nodeHasPeersForHash](10*scalingFactor, 20*scalingFactor),
						requestMetaInfo:              concurrency.NewBufferedConcurrentChannel[infoHashWithPeers](10*scalingFactor, 40*scalingFactor),
						persistTorrents:              concurrency.NewBufferedConcurrentChannel[infoHashWithMetaInfo](1000, 1),
						saveTorrentsRoot:             params.Config.SaveTorrentsRoot,
						saveTorrentsTempSuffix:       params.Config.SaveTorrentsTempSuffix,
						ignoreHashes: &ignoreHashes{
							bloom: boom.NewStableBloomFilter(10_000_000, 2, 0.001),
						},
						soughtNodeID: &concurrency.AtomicValue[protocol.ID]{},
						stopped:      make(chan struct{}),
						logger:       params.Logger.Named("dht_crawler"),
					}
					if absRoot, err := filepath.Abs(c.saveTorrentsRoot); err != nil {
						c.logger.Warnw("failed to resolve torrents path, using original", "path", c.saveTorrentsRoot, "error", err)
					} else {
						c.saveTorrentsRoot = absRoot
					}
					initialSoughtID := protocol.RandomNodeID()
					c.soughtNodeID.Set(initialSoughtID)
					c.logger.Infow(
						"crawler starting",
						"nodeID", params.KTable.Origin().String(),
						"initialSoughtNodeID", initialSoughtID.String(),
						"bootstrapNodes", len(params.Config.BootstrapNodes),
						"saveTorrentsRoot", c.saveTorrentsRoot,
					)
					go c.start()
					return nil
				},
				OnStop: func(context.Context) error {
					if c.stopped != nil {
						close(c.stopped)
					}
					return nil
				},
			},
		),
	}
}
