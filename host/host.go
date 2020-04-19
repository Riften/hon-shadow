package host

import (
	"context"
	"github.com/libp2p/go-libp2p"
	p2phost "github.com/libp2p/go-libp2p-core/host"
)

// Create host for shadow peer

func NewHost(ctx context.Context) (p2phost.Host, error){
	return libp2p.New(ctx)
}

func NewHost2(ctx context.Context)(p2phost.Host, error){
	return nil, nil
}