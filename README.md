# bitmagnet

A self-hosted BitTorrent indexer, DHT crawler, content classifier and torrent search engine with web UI, GraphQL API and Servarr stack integration.

Visit the website at [bitmagnet.io](https://bitmagnet.io).

## Runtime configuration

The minimal crawler exposes its settings through environment variables. Values are read at startup; any field not provided falls back to the compiled default. Units follow Go’s duration syntax (e.g. `30s`, `5m`).

| Variable | Example | Default | Description |
| --- | --- | --- | --- |
| `DHT_SERVER_PORT` | `6881` | `3334` | UDP port the embedded DHT node listens on. |
| `DHT_SERVER_QUERY_TIMEOUT` | `4s` | `4s` | Timeout for DHT RPC queries (ping, find_node, get_peers, etc.). |
| `DHT_CRAWLER_SCALING_FACTOR` | `8` | `10` | Multiplier that sizes internal queues and worker concurrency. |
| `DHT_CRAWLER_BOOTSTRAP_NODES` | `router.utorrent.com:6881,router.bittorrent.com:6881` | `router.utorrent.com:6881, router.bittorrent.com:6881, dht.transmissionbt.com:6881, dht.aelitis.com:6881, router.silotis.us:6881, dht.libtorrent.org:25401` | Comma-separated list of DHT routers used to seed the routing table. |
| `DHT_CRAWLER_RESEED_BOOTSTRAP_NODES_INTERVAL` | `10m` | `1m` | How often the crawler re-contacts bootstrap routers. |
| `DHT_CRAWLER_SAVE_TORRENTS_ROOT` | `/data/torrents` | `./torrents` | Directory where fetched `.tfile` bundles are written. |
| `METAINFO_REQUESTER_REQUEST_TIMEOUT` | `8s` | `6s` | Timeout for the BitTorrent metadata handshake/transfer. |
| `METAINFO_REQUESTER_KEY_MUTEX_SIZE` | `2048` | `1000` | Size of the keyed limiter that throttles concurrent metadata requests per infohash. |
