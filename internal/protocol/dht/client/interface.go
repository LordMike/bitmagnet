package client

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"net/netip"
)

type Client interface {
	Ping(ctx context.Context, addr netip.AddrPort) (PingResult, error)
	FindNode(ctx context.Context, addr netip.AddrPort, target protocol.ID) (FindNodeResult, error)
	GetPeers(ctx context.Context, addr netip.AddrPort, infoHash protocol.ID) (GetPeersResult, error)
	SampleInfoHashes(ctx context.Context, addr netip.AddrPort, target protocol.ID) (SampleInfoHashesResult, error)
}

type PingResult struct {
	ID protocol.ID
}

type FindNodeResult struct {
	ID    protocol.ID
	Nodes []NodeInfo
}

type GetPeersResult struct {
	ID     protocol.ID
	Values []netip.AddrPort
	Nodes  []NodeInfo
}

type SampleInfoHashesResult struct {
	ID       protocol.ID
	Samples  []protocol.ID
	Nodes    []NodeInfo
	Num      int
	Interval int
}

type NodeInfo struct {
	ID   protocol.ID
	Addr netip.AddrPort
}
