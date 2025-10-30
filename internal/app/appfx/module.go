package appfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/app/boilerplateappfx"
	"github.com/bitmagnet-io/bitmagnet/internal/dhtcrawler/dhtcrawlerfx"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/dhtfx"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainfofx"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"app",
		boilerplateappfx.New(),
		dhtfx.New(),
		metainfofx.New(),
		dhtcrawlerfx.New(),
	)
}
