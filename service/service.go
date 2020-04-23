package service

import (
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"time"
)

var dhtReadMessageTimeout = time.Minute
var dhtStreamIdleTimeout = 10 * time.Minute
var ErrReadTimeout = fmt.Errorf("timed out reading response")

type Service interface{
	Start()			// Register service to host.
	Protocol() protocol.ID	// Provide protocal for sender
	GetSender(p peer.ID) (*messageSender, bool)
	AddSender(ms *messageSender)
	RemoveSender(p peer.ID)
}
